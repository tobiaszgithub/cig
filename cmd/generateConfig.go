/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "generate-config",
	Short: "Generate config file",
	Long: `Generate configuration file. This file is nessesary for the operation
of the cig tool. Configuration file should be placed in working directory or userhome/.cig/ directory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generateConfig called")

	},
}

func init() {
	rootCmd.AddCommand(generateConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
