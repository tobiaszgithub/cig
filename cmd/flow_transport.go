/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// flowTransportCmd represents the flowTransport command
var flowTransportCmd = &cobra.Command{
	Use:   "transport [source-flow-id] [destination-flow-id]",
	Short: "Transport an integration flow between systems",
	Long: `You can use the following subcommand to transport
an integration flow of designtime between systems. `,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("flowTransport called")
		if len(args) == 0 {
			log.Fatal("Required parameter source-flow-id not set")
		}
		if len(args) == 1 {
			log.Fatal("Required parameter destination-flow-id not set")
		}
		destTenantKey, _ := cmd.Flags().GetString("dest-tenant-key")
		if destTenantKey == "" {
			log.Fatal("Required flag dest-tenant-key not set")
		}
		destFlowName, _ := cmd.Flags().GetString("dest-flow-name")
		destPackageId, _ := cmd.Flags().GetString("dest-package-id")

		client.RunTransportFlow(args[0], args[1], destTenantKey, destFlowName, destPackageId)
	},
}

func init() {
	flowCmd.AddCommand(flowTransportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowTransportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowTransportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flowTransportCmd.Flags().StringP("dest-tenant-key", "t", "", "Destination tenant key from configuration file")
	flowTransportCmd.Flags().StringP("dest-flow-name", "n", "", "Destination Integration Flow name")
	flowTransportCmd.Flags().StringP("dest-package-id", "p", "", "Destination Integration Flow package id")
}
