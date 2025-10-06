package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/goloop/env"
)

type errEnvVariableNotFound struct {
	Variables []error
}

func (e *errEnvVariableNotFound) Error() string {
	if len(e.Variables) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("missing or invalid configuration:\n")
	for _, err := range e.Variables {
		sb.WriteString(" - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *errEnvVariableNotFound) Add(err error) {
	e.Variables = append(e.Variables, err)
}

// Config struct
type Config struct {
	DB             *DB
	Closer         *Closer
	HttpServer     *HttpServer
	GrpcAuthServer *GRPCAuthServer
}

func (c *Config) Close(ctx context.Context) error {
	const op = "config.Close"
	env.Clear()
	return nil
}

// New load parameters from the env file and return Config
func New() (*Config, error) {
	err := MustLoad(".env")
	if err != nil {
		return nil, err
	}
	var cfg = &Config{}

	var envErrs = &errEnvVariableNotFound{}

	cfg.LoadDbConfig(envErrs)
	cfg.LoadCloserConfig(envErrs)
	cfg.LoadHttpServerConfig(envErrs)
	cfg.LoadGRPCConfig(envErrs)

	if envErrs.Variables != nil {
		if len(envErrs.Variables) > 0 {
			return nil, envErrs
		}
	}

	return cfg, nil
}

// MustLoad loading parameters from the env file
func MustLoad(filename string) error {
	const op = "config.MustLoad"
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("op: %s, file: %s does not exist", op, filename)
		} else {
			return err
		}
	}

	err = env.Load(filename)
	if err != nil {
		return fmt.Errorf("op: %s,incorrect data in the configuration file: %s", op, err.Error())
	}
	return nil
}
