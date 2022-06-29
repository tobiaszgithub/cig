/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// flowConfigurationsCmd represents the flowConfigurations command
var flowConfigurationsCmd = &cobra.Command{
	Use:     "describe-configs [flow id]",
	Aliases: []string{"configs", "configurations"},
	Short:   "Get configurations of an integration flow by Id and version",
	Long: `You can use the following request to get the configuration
parameters (key/value pairs) of a designtime integration artifact by Id and version.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.RunGetFlowConfigs(args[0])
	},
}

func init() {
	flowCmd.AddCommand(flowConfigurationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowConfigurationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowConfigurationsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
