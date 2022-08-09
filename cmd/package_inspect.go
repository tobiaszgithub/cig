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

// packageInspectCmd represents the packageInspect command
var packageInspectCmd = &cobra.Command{
	Use:   "inspect package-id",
	Short: "Get integration package by ID",
	Long:  `You can use the following subcommand to get an integration packages of designtime by Id.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) == 0 {
			log.Fatal("Required parameter package-id not set")
		}
		client.RunInspectIntegrationPackage(conf, args[0])
	},
}

func init() {
	packageCmd.AddCommand(packageInspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageInspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageInspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
