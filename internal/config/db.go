package config

import (
	"fmt"
	"os"
	"time"
)

//type DBConfig interface {
//	DSN() string
//}

type DB struct {
	dsn            string
	connectTimeout time.Duration
}

func (c *Config) LoadDbConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadDbConfig"
	cfg := &DB{}

	if env, ok := os.LookupEnv("DB_URL"); ok {
		cfg.dsn = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'DB_URL' is not set", op))
	}

	if env, ok := os.LookupEnv("DB_CONNECT_TIMEOUT"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'DB_CONNECT_TIMEOUT' bad format", op))
		}
		cfg.connectTimeout = t
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'DB_CONNECT_TIMEOUT' is not set", op))
	}

	c.DB = cfg
}

func (cfg *DB) DSN() string {
	return cfg.dsn
}

func (cfg *DB) ConnectTimeout() time.Duration {
	return cfg.connectTimeout
}
