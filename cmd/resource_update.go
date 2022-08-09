/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
	"github.com/tobiaszgithub/cig/config"
)

// resourceUpdateCmd represents the resourceUpdate command
var resourceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a resource of an integration flow",
	Long:  `You can use the following command to update a resource of an integration flow from designtime.`,
	Run: func(cmd *cobra.Command, args []string) {

		conf, err := config.NewConfiguration(TenantKey)
		if err != nil {
			log.Fatal(err)
		}

		flowId, _ := cmd.Flags().GetString("flow-id")
		flowVersion, _ := cmd.Flags().GetString("flow-version")
		resourceName := args[0]
		resourceType, _ := cmd.Flags().GetString("resource-type")
		resourceFileName, _ := cmd.Flags().GetString("resource-file-name")

		client.RunResourceUpdate(os.Stdout, conf, flowId, flowVersion, resourceName, resourceType, resourceFileName)
	},
}

func init() {
	resourceCmd.AddCommand(resourceUpdateCmd)

	resourceUpdateCmd.Flags().StringP("flow-id", "i", "", "Integration Flow id")
	resourceUpdateCmd.Flags().StringP("flow-version", "v", "active", "Integration Flow version")
	resourceUpdateCmd.Flags().StringP("resource-type", "y", "groovy", "Resource type. Available values: edmx, groovy, jar, js, mmap, opmap, wsdl, xsd, xslt")
	resourceUpdateCmd.Flags().StringP("resource-file-name", "f", "", "Resource file (.groovy,.js,.wsdl...)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resourceUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resourceUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
