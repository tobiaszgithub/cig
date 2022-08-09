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

// flowCopyCmd represents the flowCopy command
var flowCopyCmd = &cobra.Command{
	Use:   "copy [source-flow-id] [destination-flow-id]",
	Short: "Copy an integration flow",
	Long: `You can use the following subcommand to copy
an integration flow of designtime. `,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) == 0 {
			log.Fatal("Required parameter source-flow-id not set")
		}
		if len(args) == 1 {
			log.Fatal("Required parameter destination-flow-id not set")
		}
		destFlowName, _ := cmd.Flags().GetString("dest-flow-name")
		destPackageId, _ := cmd.Flags().GetString("dest-package-id")

		client.RunCopyFlow(conf, args[0], args[1], destFlowName, destPackageId)

	},
}

func init() {
	flowCmd.AddCommand(flowCopyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowCopyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowCopyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flowCopyCmd.Flags().StringP("dest-flow-name", "n", "", "Destination Integration Flow name")
	flowCopyCmd.Flags().StringP("dest-package-id", "p", "", "Destination Integration Flow package id")
}
