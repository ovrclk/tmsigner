package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagJSON = "json"
)

func jsonFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagJSON, "j", false, "returns the response in json format")
	if err := viper.BindPFlag(flagJSON, cmd.Flags().Lookup(flagJSON)); err != nil {
		panic(err)
	}
	return cmd
}
