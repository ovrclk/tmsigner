/*
Copyright © 2020 Jack Zampolin jack.zampolin@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	cfgPath     string
	homePath    string
	debug       bool
	config      *Config
	defaultHome = os.ExpandEnv("$HOME/.tmsigner")
)

func init() {
	cobra.EnableCommandSorting = false

	rootCmd.SilenceUsage = true

	// Register top level flags --home and --config
	// TODO: just rely on homePath and remove the config path arg?
	rootCmd.PersistentFlags().StringVar(&homePath, flags.FlagHome, defaultHome, "set home directory")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	if err := viper.BindPFlag(flags.FlagHome, rootCmd.Flags().Lookup(flags.FlagHome)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug")); err != nil {
		panic(err)
	}

	// Register subcommands
	rootCmd.AddCommand(
		startCmd(),
		configInitCmd(),
	)

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmsigner",
	Short: "This application signs blocks for tendermint using a configured private key",
	Long:  strings.TrimSpace(``),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		// reads `homeDir/config/config.yaml` into `var config *Config` before each command
		return initConfig(rootCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}

	config = &Config{}
	cfgPath := path.Join(home, "config.yaml")
	if _, err := os.Stat(cfgPath); err == nil {
		viper.SetConfigFile(cfgPath)
		if err := viper.ReadInConfig(); err == nil {
			config, err = LoadConfigFromFile(viper.ConfigFileUsed())
			if err != nil {
				fmt.Println("Error reading in config:", err)
				os.Exit(1)
			}
		}
	}
	return nil
}

func overWriteConfig(cmd *cobra.Command, cfg *Config) error {
	home, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}

	cfgPath := path.Join(home, "config.yaml")
	if _, err = os.Stat(cfgPath); err == nil {
		viper.SetConfigFile(cfgPath)
		if err = viper.ReadInConfig(); err == nil {

			// marshal the new config
			out, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}

			// overwrite the config file
			err = ioutil.WriteFile(viper.ConfigFileUsed(), out, 0600)
			if err != nil {
				return err
			}

			// set the global variable
			config = cfg
		}
	}
	return err
}

// readLineFromBuf reads one line from stdin.
func readStdin() (string, error) {
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(str), err
}
