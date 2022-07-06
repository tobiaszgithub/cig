/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// flowDownloadCmd represents the flowDownload command
var flowDownloadCmd = &cobra.Command{
	Use:   "download [flow id]",
	Short: "Download an integration flow as zip file",
	Long: `You can use the following subcommand to download an integration flow of designtime as zip file.
Integration flows of configure-only packages cannot be downloaded.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.RunDownloadFlow(args[0])
	},
}

func init() {
	flowCmd.AddCommand(flowDownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flowDownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flowDownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
