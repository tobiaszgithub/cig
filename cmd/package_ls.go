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

// packageLsCmd represents the packageLs command
var packageLsCmd = &cobra.Command{
	Use:   "ls [package-id]",
	Short: "Get all integration packages as list or get all integration flow of the package",
	Long: `You can use the following subcommand to get all integration packages of designtime.
Optionaly you can use this subcommand to get all integration flows of the specified package-id`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("packageLs called")
		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) > 0 {
			client.RunGetFlowsOfIntegrationPackage(conf, args[0])
		} else {
			client.RunGetIntegrationPackages(conf)
		}

	},
}

func init() {
	packageCmd.AddCommand(packageLsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageLsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageLsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
