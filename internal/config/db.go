package config

import (
	"fmt"
	"os"
)

//type DBConfig interface {
//	DSN() string
//}

type DB struct {
	dsn string
}

func (c *Config) LoadDbConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadDbConfig"
	cfg := &DB{}

	if value, ok := os.LookupEnv("DB_URL"); ok {
		cfg.dsn = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'DB_URL' is not set", op))
	}

	c.DB = cfg
}

func (cfg *DB) DSN() string {
	return cfg.dsn
}
