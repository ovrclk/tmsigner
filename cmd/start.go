/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"net"
	"os"
	"sync"
	"time"

	"github.com/ovrclk/tmsigner/signer"

	tos "github.com/tendermint/tendermint/libs/os"
	svc "github.com/tendermint/tendermint/libs/service"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			logger := config.Logger()
			logger.Info("Tendermint Validator", "priv-key", config.PrivValKeyFile, "priv-state-dir", config.PrivValStateDir)
			signer.InitSerialization()

			// services to stop on shutdown
			var services []svc.Service

			if err = config.PrivValStateExists(); err != nil {
				return err
			}

			for _, node := range config.Nodes {
				signer := signer.NewNodeClient(
					node.Address,
					logger,
					config.ChainID,
					config.LoadPrivVal(),
					net.Dialer{Timeout: 30 * time.Second},
				)

				if err := signer.Start(); err != nil {
					panic(err)
				}

				services = append(services, signer)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			tos.TrapSignal(logger, func() {
				for _, service := range services {
					err := service.Stop()
					if err != nil {
						panic(err)
					}
				}
				wg.Done()
			})
			wg.Wait()
			return nil
		},
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
