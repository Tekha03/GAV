package config

type GRPCConfig struct {
	Addr string
}

func loadGRPC() GRPCConfig {
	return GRPCConfig{Addr: getEnv("GRPC_ADDR", ":9000")}
}
