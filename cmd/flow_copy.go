/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// flowCopyCmd represents the flowCopy command
var flowCopyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy an integration flow",
	Long: `You can use the following subcommand to copy
an integration flow of designtime`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flowCopy called")
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
}
