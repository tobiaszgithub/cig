/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/tobiaszgithub/cig/client"
	"github.com/tobiaszgithub/cig/model"
)

// updateConfigsCmd represents the updateConfigs command
var updateConfigsCmd = &cobra.Command{
	Use:   "update-configs flow-id",
	Short: "Update configuration parameters of an integration flow",
	Long: `You can use the following command to update the value
for a configuration parameters of a designtime integration flow.`,
	Run: func(cmd *cobra.Command, args []string) {

		var allConfigParams []model.FlowConfigurationPrinter

		parameters, _ := cmd.Flags().GetStringArray("parameter")
		fileWithConfigsName, _ := cmd.Flags().GetString("input-file")
		var decodedFile model.FlowConfigurationsPrinter

		if fileWithConfigsName != "" {
			fileWithConfigs, err := os.Open(fileWithConfigsName)
			if err != nil {
				log.Fatal("Error reading file: ", err)
			}
			defer fileWithConfigs.Close()

			if err := json.NewDecoder(fileWithConfigs).Decode(&decodedFile); err != nil {
				log.Fatal("Error decodeing file: ", err)
			}
		}

		configParams, _ := parseConfigParameters(parameters)

		allConfigParams = append(decodedFile.D.Results, configParams...)

		client.RunUpdateFlowConfigs(args[0], allConfigParams)

	},
}

func parseConfigParameters(parameters []string) ([]model.FlowConfigurationPrinter, error) {
	var configs []model.FlowConfigurationPrinter

	for _, p := range parameters {

		key, value, err := parseParameter(p)
		if err != nil {
			return nil, err
		}

		configParam := model.FlowConfigurationPrinter{
			ParameterKey:   key,
			ParameterValue: value,
			DataType:       "",
		}

		configs = append(configs, configParam)
	}

	return configs, nil

}

func parseParameter(param string) (string, string, error) {
	//example: Key=key1,Value=value1
	reg := regexp.MustCompile(`Key=.*,Value=`)
	key := reg.FindString(param)
	key = key[4 : len(key)-7]

	reg = regexp.MustCompile(`,Value=.*`)
	value := reg.FindString(param)
	value = value[7:]

	return key, value, nil
}

func init() {
	flowCmd.AddCommand(updateConfigsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateConfigsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateConfigsCmd.Flags().StringArrayP("parameter", "p", []string{}, "Flow Configuration parameter, format: Key=key1,Value=value1")
	updateConfigsCmd.Flags().StringP("input-file", "f", "", "File with parameters, utf-8 file has format like output from describe-configs")
	//	updateConfigsCmd.Flags().StringSliceP("parameters2", "r", []string{}, "Help message for toggle")
}
