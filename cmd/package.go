/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"ls", "p"},
	Short:   "Command related to the processing of integration packages",
	Long:    `Command related to the processing of integration packages`,
	Run: func(cmd *cobra.Command, args []string) {
		//	fmt.Println("package called")
		//cmd.packageLsCmd()
		//client.RunGetIntegrationPackages()

		if len(args) > 0 {
			client.RunGetFlowsOfIntegrationPackage(args[0])
		} else {
			client.RunGetIntegrationPackages()
		}
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
