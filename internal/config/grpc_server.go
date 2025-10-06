package config

import (
	"fmt"
	"os"
)

//type GRPCServerConfig interface {
//	Addr() string
//	Port() string
//	Host() string
//	Network() string
//}

type GRPCAuthServer struct {
	host    string
	port    string
	network string
}

func (c *Config) LoadGRPCConfig(fErr *errEnvVariableNotFound) {
	const op = "Config.LoadGRPCConfig"
	cfg := &GRPCAuthServer{}

	if value, ok := os.LookupEnv("GRPC_HOST"); ok {
		cfg.host = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_HOST' is not set", op))
	}

	if value, ok := os.LookupEnv("GRPC_PORT"); ok {
		cfg.port = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_PORT' is not set", op))
	}

	if value, ok := os.LookupEnv("GRPC_NETWORK"); ok {
		cfg.network = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_NETWORK' is not set", op))
	}

	c.GrpcAuthServer = cfg
}

func (c *GRPCAuthServer) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *GRPCAuthServer) Port() string {
	return c.port
}

func (c *GRPCAuthServer) Host() string {
	return c.host
}

func (c *GRPCAuthServer) Network() string {
	return c.network
}
