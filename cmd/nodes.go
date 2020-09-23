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
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

// chainCmd represents the keys command
func nodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nodes",
		Aliases: []string{"n"},
		Short:   "manage nodes that the signer connects to",
	}

	cmd.AddCommand(nodesList())
	cmd.AddCommand(nodesDelete())
	cmd.AddCommand(nodesAdd())

	return cmd
}

func nodesList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list configured nodes",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			for i, n := range config.Nodes {
				fmt.Printf("node %d: %s\n", i, n.Address)
			}
			return nil
		},
	}
	return cmd
}

func nodesAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [address]",
		Aliases: []string{"a"},
		Short:   "add node to list",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if _, err = url.ParseRequestURI(args[0]); err != nil {
				return err
			}
			c := config
			c.Nodes = append(c.Nodes, &NodeConfig{Address: args[0]})
			return overWriteConfig(c)
		},
	}
	return cmd
}

func nodesDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [index]",
		Aliases: []string{"d"},
		Short:   "delete node from list, check 'nodes list' to get the index of the node you wish to remove",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			index, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			c := &Config{ChainID: config.ChainID, Nodes: []*NodeConfig{}, home: config.home}
			for i, n := range config.Nodes {
				if i != int(index) {
					c.Nodes = append(c.Nodes, n)
				}
			}

			return overWriteConfig(c)
		},
	}
	return cmd
}
