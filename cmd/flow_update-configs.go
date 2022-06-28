/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
)

// updateConfigsCmd represents the updateConfigs command
var updateConfigsCmd = &cobra.Command{
	Use:   "update-configs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		parameters, _ := cmd.Flags().GetStringArray("parameter")

		fmt.Println(parameters)

		fmt.Println("updateConfigs called")

		client.RunUpdateFlowConfigs(args[0], parameters)

		fmt.Println(args)
	},
}

func init() {
	flowCmd.AddCommand(updateConfigsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateConfigsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateConfigsCmd.Flags().StringArrayP("parameter", "p", []string{}, "Help message for toggle")
	//	updateConfigsCmd.Flags().StringSliceP("parameters2", "r", []string{}, "Help message for toggle")
}
