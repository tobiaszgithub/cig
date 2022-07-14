/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/config"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "generate-config",
	Short: "Generate config file",
	Long: `Generate configuration file. This file is nessesary for the operation
of the cig tool. Configuration file should be placed in working directory or userhome/.cig/ directory`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("generateConfig called")
		outputFileName, _ := cmd.Flags().GetString("output-file")
		log.Println("File: ", outputFileName, " will be generated")
		err := config.GenerateEmptyConfigFile(outputFileName)
		if err != nil {
			log.Fatal("Error during generating Configuration file: ", err)
		}
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
	generateConfigCmd.Flags().StringP("output-file", "o", "config.json", "The output file with empty configuration parameters that will be created")
	//generateConfigCmd.MarkFlagRequired("output-file")
}
