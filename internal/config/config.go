package config

import "github.com/brocaar/chirpstack-packet-multiplexer/internal/multiplexer"

// Version holds the ChirpStack Packet Multiplexer version.
var Version string

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel int `mapstructure:"log_level"`
	} `mapstructure:"general"`

	PacketMultiplexer multiplexer.Config `mapstructure:"packet_multiplexer"`
}

// C holds the configuration.
var C Config
