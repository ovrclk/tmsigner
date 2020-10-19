package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/simapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	tmconfig "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Command for inititalizing an empty config at the --home location
func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [chain-id]",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a default home directory at path defined by --home",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := defaultHome
			dataDir := path.Join(home, "data")
			cfgDir := path.Join(home, "config")
			cfgPath := path.Join(home, "config.toml")

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
					// Create the config folder in the home dir to
					// ensure proper placement of node key
					if err = os.Mkdir(cfgDir, os.ModePerm); err != nil {
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
	ChainID string        `toml:"chain_id"`
	Nodes   []*NodeConfig `toml:"node"`

	home string
}

// TMConfig returns the tendermint configuration
func (c *Config) TMConfig() *tmconfig.Config {
	c.home = defaultHome
	return &tmconfig.Config{
		BaseConfig: tmconfig.BaseConfig{
			RootDir:            c.home,
			Moniker:            c.ChainID,
			Genesis:            c.GenesisFile(),
			PrivValidatorKey:   c.PrivValKeyFile(),
			PrivValidatorState: c.PrivValStateFile(),
			NodeKey:            c.NodeKeyFile(),
		},
		RPC:             tmconfig.DefaultRPCConfig(),
		P2P:             tmconfig.DefaultP2PConfig(),
		Mempool:         tmconfig.DefaultMempoolConfig(),
		StateSync:       tmconfig.DefaultStateSyncConfig(),
		FastSync:        tmconfig.DefaultFastSyncConfig(),
		Consensus:       tmconfig.DefaultConsensusConfig(),
		TxIndex:         tmconfig.DefaultTxIndexConfig(),
		Instrumentation: tmconfig.DefaultInstrumentationConfig(),
	}
}

// CLIContext returns the proper cliContext for the config
func (c *Config) CLIContext() client.Context {
	encodingConfig := simapp.MakeEncodingConfig()
	return client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(c.home)
}

// GenesisFile returns the location of the GenesisFile
func (c *Config) GenesisFile() string {
	return filepath.Join(c.home, "genesis.json")
}

// PrivValKeyFile returns the location of the PrivValKeyFile
func (c *Config) PrivValKeyFile() string {
	return filepath.Join(c.home, "priv_validator_key.json")
}

// PrivValStateDir returns the location of the PrivValStateDir
func (c *Config) PrivValStateDir() string {
	return filepath.Join(c.home, "data")
}

// NodeKeyFile used to initialize a node key for signing gentx
func (c *Config) NodeKeyFile() string {
	return filepath.Join(c.home, "config", "node_key.json")
}

// Logger returns the tendermint logger
func (c *Config) Logger() tmlog.Logger {
	return tmlog.NewTMLogger(
		tmlog.NewSyncWriter(os.Stdout),
	).With("module", "validator")
}

// PrivValStateFile returns the path to the priv_validator_state.json file for the instance
func (c *Config) PrivValStateFile() string {
	return path.Join(c.PrivValStateDir(), "priv_validator_state.json")
}

// PrivValStateExists returns an error if the priv val state doesn't exist
func (c *Config) PrivValStateExists() error {
	if !fileExists(c.PrivValStateFile()) {
		return fmt.Errorf("state file missing: %s", c.PrivValStateFile())
	}
	return nil
}

// LoadPrivVal returns the parsed priv validator json
func (c *Config) LoadPrivVal() tmtypes.PrivValidator {
	return privval.LoadFilePV(c.PrivValKeyFile(), c.PrivValStateFile())
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

	config.home = filepath.Dir(file)

	return config, nil
}

func overWriteConfig(cfg *Config) (err error) {
	cfgPath := path.Join(cfg.home, "config.toml")
	if _, err = os.Stat(cfgPath); err == nil {
		buff := bytes.NewBuffer([]byte{})
		if err := toml.NewEncoder(buff).Encode(cfg); err != nil {
			panic(err)
		}

		// overwrite the config file
		err = ioutil.WriteFile(cfgPath, buff.Bytes(), 0600)
		if err != nil {
			return err
		}

		// set the global variable
		config = cfg
	}
	return err
}
