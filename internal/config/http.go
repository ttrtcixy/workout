package config

import (
	"fmt"
	"os"
	"time"
)

type HttpServer struct {
	host            string
	port            string
	shutdownTimeout time.Duration
}

func (c *Config) LoadHttpServerConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadSmtpConfig"
	cfg := &HttpServer{}

	if value, ok := os.LookupEnv("HTTP_HOST"); ok {
		cfg.host = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'HTTP_HOST' is not set", op))
	}

	if value, ok := os.LookupEnv("HTTP_PORT"); ok {
		cfg.port = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'HTTP_PORT' is not set", op))
	}

	if value, ok := os.LookupEnv("HTTP_SHUTDOWN_TIME"); ok {
		dur, err := time.ParseDuration(value)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'HTTP_SHUTDOWN_TIME' bad format", op))
		}
		cfg.shutdownTimeout = dur
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'HTTP_SHUTDOWN_TIME' is not set", op))
	}

	c.HttpServer = cfg
}

func (c *HttpServer) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *HttpServer) ShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}
