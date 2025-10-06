package config

import (
	"fmt"
	"os"
	"time"
)

type CloserConfig struct {
	totalDuration time.Duration
	funcDuration  time.Duration
}

func (c *CloserConfig) TotalDuration() time.Duration {
	return c.totalDuration
}

func (c *CloserConfig) FuncDuration() time.Duration {
	return c.funcDuration
}

func (c *Config) LoadCloserConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadCloserConfig"

	var cfg = &CloserConfig{}
	if env, ok := os.LookupEnv("CLOSER_TOTAL_DURATION"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'CLOSER_TOTAL_DURATION' bad format", op))
		}
		cfg.totalDuration = t
	} else {
		cfg.totalDuration = 0
		//fErr.Add(fmt.Errorf("%s: env variable 'CLOSER_TOTAL_DURATION' is not set", op))
	}

	if env, ok := os.LookupEnv("CLOSER_FUNC_DURATION"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'CLOSER_FUNC_DURATION' bad format", op))
		}
		cfg.funcDuration = t
	} else {
		cfg.totalDuration = 0
		//fErr.Add(fmt.Errorf("%s: env variable 'CLOSER_FUNC_DURATION' is not set", op))
	}

	c.CloserConfig = cfg
}
