package signer

import (
	"os"

	"github.com/BurntSushi/toml"
)

// NodeConfig contains the configuration for an individual node
type NodeConfig struct {
	Address string `toml:"address"`
}

// Config represents the configuration file
type Config struct {
	PrivValKeyFile  string `toml:"key_file"`
	PrivValStateDir string `toml:"state_dir"`
	ChainID         string `toml:"chain_id"`

	Nodes []NodeConfig `toml:"node"`
}

// LoadConfigFromFile returns the config struct from the file
func LoadConfigFromFile(file string) (*Config, error) {
	var config *Config

	reader, err := os.Open(file)
	if err != nil {
		return config, err
	}
	if _, err = toml.DecodeReader(reader, config); err != nil {
		return nil, err
	}
	return config, nil
}
