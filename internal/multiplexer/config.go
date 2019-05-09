package multiplexer

// Config holds the multiplexer config.
type Config struct {
	Bind     string          `mapstructure:"bind"`
	Backends []BackendConfig `mapstructure:"backend"`
}

// BackendConfig holds the config for a single backend.
type BackendConfig struct {
	Host       string   `mapstructure:"host"`
	UplinkOnly bool     `mapstructure:"uplink_only"`
	GatewayIDs []string `mapstructure:"gateway_ids"`
}
