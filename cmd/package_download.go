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

// packageDownloadCmd represents the packageDownload command
var packageDownloadCmd = &cobra.Command{
	Use:   "download package-id",
	Short: "Download integration package by ID",
	Long: `You can use the following subcommand to download an integration package of designtime as .zip file.
Download fails if the package contains one or more artifacts in draft state.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Fatal("Required parameter package-id not set")
		}
		client.RunDownloadIntegrationPackage(conf, args[0])

	},
}

func init() {
	packageCmd.AddCommand(packageDownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageDownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageDownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
