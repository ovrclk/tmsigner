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
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/privval"
)

// chainCmd represents the keys command
func privValCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "priv-val",
		Aliases: []string{"pv"},
		Short:   "Commands to manage the priv-validator files in the signer",
	}

	cmd.AddCommand(pvStateShow())
	cmd.AddCommand(pvCreate())

	return cmd
}

func pvStateShow() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s"},
		Short:   "Show the current round state from the database",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			pv := config.LoadPrivVal()
			out, ok := pv.(*privval.FilePV)
			if !ok {
				panic(err)
			}
			bz, err := json.Marshal(out.LastSignState)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}
func pvCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{},
		Short:   "create a priv_validator_key and priv_validator_state file USUSAL WARNINGS APPLY!",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_, _, err = genutil.InitializeNodeValidatorFiles(config.TMConfig())
			return err
		},
	}
	return cmd
}
