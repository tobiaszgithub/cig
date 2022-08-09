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

// createCmd represents the create command
var flowCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or upload an integration flow",
	Long: `You can use the following subcommand to create or upload
an integration flow of designtime`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		packageid, _ := cmd.Flags().GetString("package-id")
		fileName, _ := cmd.Flags().GetString("content-file-name")

		client.RunCreateFlow(conf, name, id, packageid, fileName)
	},
}

func init() {
	flowCmd.AddCommand(flowCreateCmd)
	flowCreateCmd.Flags().StringP("name", "n", "", "Integration Flow name")
	flowCreateCmd.Flags().StringP("id", "i", "", "Integration Flow id")
	flowCreateCmd.Flags().StringP("package-id", "p", "", "Integration Flow package id")
	flowCreateCmd.Flags().StringP("content-file-name", "f", "", "Integration Flow artifact content file (.zip)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
