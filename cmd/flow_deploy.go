/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
	"github.com/tobiaszgithub/cig/config"
)

// flowDeployCmd represents the flowDeploy command
var flowDeployCmd = &cobra.Command{
	Use:   "deploy flow-id",
	Short: "Deploy an integration flow",
	Long:  `You can use the following request to deploy an integration flow of designtime.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Fatal("Required parameter flow-id not set")
		}
		version, _ := cmd.Flags().GetString("version")
		client.RunDeployFlow(conf, args[0], version)
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
