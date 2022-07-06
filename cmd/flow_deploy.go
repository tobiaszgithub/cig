/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// flowDeployCmd represents the flowDeploy command
var flowDeployCmd = &cobra.Command{
	Use:   "deploy [flow-id]",
	Short: "Deploy an integration flow",
	Long:  `You can use the following request to deploy an integration flow of designtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flowDeploy called")
		version, _ := cmd.Flags().GetString("version")
		client.RunDeployFlow(args[0], version)
	},
}

func init() {
	flowCmd.AddCommand(flowDeployCmd)
	flowDeployCmd.Flags().StringP("version", "v", "active", "Integration Flow version")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowDeployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowDeployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
