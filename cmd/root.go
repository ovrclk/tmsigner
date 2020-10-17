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
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgPath     string
	homePath    string
	config      *Config
	defaultHome = os.ExpandEnv("$HOME/.tmsigner")
)

func init() {
	config = &Config{home: defaultHome}
	cobra.EnableCommandSorting = false
	rootCmd.SilenceUsage = true
	encodingConfig := simapp.MakeEncodingConfig()
	// Register subcommands

	rootCmd.AddCommand(
		startCmd(),
		configInitCmd(),
		getVersionCmd(),
		nodesCmd(),
		privValCmd(),

		keys.Commands(defaultHome),
		txCmd(encodingConfig),
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
		// reads `homeDir/config.toml` into `var config *Config` before each command
		return initConfig(rootCmd)
	}

	ctx := context.Background()
	cliCtx := client.Context{}.WithJSONMarshaler(simapp.MakeEncodingConfig().Marshaler)
	ctx = context.WithValue(ctx, client.ClientContextKey, &cliCtx)
	ctx = context.WithValue(ctx, server.ServerContextKey, &server.Context{
		Viper:  viper.New(),
		Config: config.TMConfig(),
		Logger: config.Logger(),
	})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	config = &Config{home: defaultHome}
	cfgPath := path.Join(defaultHome, "config.toml")
	if _, err := os.Stat(cfgPath); err == nil {
		config, err = LoadConfigFromFile(cfgPath)
		if err != nil {
			fmt.Println("Error reading in config:", err)
			os.Exit(1)
		}
	}

	return nil
}

// readLineFromBuf reads one line from stdin.
func readStdin() (string, error) {
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(str), err
}
