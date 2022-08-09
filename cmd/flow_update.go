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

// updateFlowCmd represents the updateFlow command
var flowUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an integration flow",
	Long:  `You can use the following request to update an integration flow from designtime`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		fileName, _ := cmd.Flags().GetString("content-file-name")
		version, _ := cmd.Flags().GetString("version")

		client.RunUpdateFlow(conf, name, id, version, fileName)
	},
}

func init() {
	flowCmd.AddCommand(flowUpdateCmd)

	flowUpdateCmd.Flags().StringP("name", "n", "", "Integration Flow name")
	flowUpdateCmd.Flags().StringP("id", "i", "", "Integration Flow id")
	flowUpdateCmd.Flags().StringP("version", "v", "active", "Integration Flow version")
	flowUpdateCmd.Flags().StringP("content-file-name", "f", "", "Integration Flow artifact content file (.zip)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateFlowCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateFlowCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
