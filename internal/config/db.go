package config

import (
	"fmt"
	"os"
)

//type DBConfig interface {
//	DSN() string
//}

type DBConfig struct {
	dsn string
}

func (c *Config) LoadDbConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadDbConfig"
	cfg := &DBConfig{}

	if value, ok := os.LookupEnv("DB_URL"); ok {
		cfg.dsn = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'DB_URL' is not set", op))
	}

	c.DBConfig = cfg
}

func (cfg *DBConfig) DSN() string {
	return cfg.dsn
}
