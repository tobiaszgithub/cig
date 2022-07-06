/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// flowInspectCmd represents the flowInspect command
var flowInspectCmd = &cobra.Command{
	Use:   "inspect flow-id",
	Short: "Get integration flow by id and version",
	Long:  `You can use the following subcommand to get an integration flow of designtime by Id and version.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.RunInspectFlow(args[0])
	},
}

func init() {
	flowCmd.AddCommand(flowInspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowInspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowInspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
