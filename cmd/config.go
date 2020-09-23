package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ovrclk/tmsigner/signer"
	"github.com/spf13/cobra"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
)

// Command for inititalizing an empty config at the --home location
func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [chain-id]",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a default home directory at path defined by --home",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flags.FlagHome)
			if err != nil {
				return err
			}

			dataDir := path.Join(home, "data")
			cfgPath := path.Join(home, "config.yaml")

			// If the config doesn't exist...
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				// And the config folder doesn't exist...
				if _, err := os.Stat(dataDir); os.IsNotExist(err) {
					// And the home folder doesn't exist
					if _, err := os.Stat(home); os.IsNotExist(err) {
						// Create the home folder
						if err = os.Mkdir(home, os.ModePerm); err != nil {
							return err
						}
					}
					// Create the home config folder
					if err = os.Mkdir(dataDir, os.ModePerm); err != nil {
						return err
					}
				}

				// Then create the file...
				f, err := os.Create(cfgPath)
				if err != nil {
					return err
				}
				defer f.Close()

				// And write the default config to that location...
				if _, err = f.Write(defaultConfig(args[0])); err != nil {
					return err
				}

				// And return no error...
				return nil
			}

			// Otherwise, the config file exists, and an error is returned...
			return fmt.Errorf("config already exists: %s", cfgPath)
		},
	}
	return cmd
}

// NodeConfig contains the configuration for an individual node
type NodeConfig struct {
	Address string `toml:"address"`
}

// Config represents the configuration file
type Config struct {
	PrivValKeyFile  string        `toml:"priv_val_key_file,omitempty"`
	PrivValStateDir string        `toml:"priv_val_state_file,omitempty"`
	ChainID         string        `toml:"chain_id"`
	Nodes           []*NodeConfig `toml:"node"`
}

// Logger returns the tendermint logger
func (c *Config) Logger() tmlog.Logger {
	return tmlog.NewTMLogger(
		tmlog.NewSyncWriter(os.Stdout),
	).With("module", "validator")
}

// PrivValStateFile returns the path to the priv_validator_state.json file for the instance
func (c *Config) PrivValStateFile() string {
	return path.Join(c.PrivValStateDir, fmt.Sprintf("%s_priv_validator_state.json", c.ChainID))
}

// PrivValStateExists returns an error if the priv val state doesn't exist
func (c *Config) PrivValStateExists() error {
	if !fileExists(c.PrivValStateFile()) {
		return fmt.Errorf("state file missing: %s", c.PrivValStateFile())
	}
	return nil
}

// LoadPrivVal returns the parsed priv validator json
func (c *Config) LoadPrivVal() *signer.PvGuard {
	val := privval.LoadFilePV(c.PrivValKeyFile, c.PrivValStateFile())
	return &signer.PvGuard{PrivValidator: val}
}

func defaultConfig(chainID string) []byte {
	buff := bytes.NewBuffer([]byte{})

	if err := toml.NewEncoder(buff).Encode(Config{
		ChainID: chainID,
		Nodes: []*NodeConfig{
			{Address: "tcp://localhost:1234"},
		},
	}); err != nil {
		panic(err)
	}

	return buff.Bytes()
}

// LoadConfigFromFile returns the config struct from the file
func LoadConfigFromFile(file string) (*Config, error) {
	fmt.Println(file)
	reader, err := os.Open(file)
	if err != nil {
		return config, err
	}

	if _, err = toml.DecodeReader(reader, config); err != nil {
		return nil, err
	}

	if config.ChainID == "" {
		return nil, fmt.Errorf("must provide chain_id in configuration")
	}

	if len(config.Nodes) == 0 {
		return nil, fmt.Errorf("must configure nodes to sign for")
	}

	config.PrivValKeyFile = filepath.Join(filepath.Dir(file), "priv_validator_key.json")
	config.PrivValStateDir = filepath.Join(filepath.Dir(file), "data")

	return config, nil
}
