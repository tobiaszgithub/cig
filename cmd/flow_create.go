/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// createCmd represents the create command
var createFlowCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or upload an integration flow of designtime",
	Long: `You can use the following subcommand to create or upload
	an integration flow of designtime`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")
		packageid, _ := cmd.Flags().GetString("package-id")
		fileName, _ := cmd.Flags().GetString("content-file-name")

		client.RunCreateFlow(name, id, packageid, fileName)
	},
}

func init() {
	flowCmd.AddCommand(createFlowCmd)
	createFlowCmd.Flags().StringP("name", "n", "", "Integration Flow name")
	createFlowCmd.Flags().StringP("id", "i", "", "Integration Flow id")
	createFlowCmd.Flags().StringP("package-id", "p", "", "Integration Flow package id")
	createFlowCmd.Flags().StringP("content-file-name", "f", "", "Integration Flow artifact content file (.zip)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
