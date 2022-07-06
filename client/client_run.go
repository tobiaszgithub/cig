package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

func RunGetIntegrationPackages() {

	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := GetIntegrationPackages(conf)

	if err != nil {
		log.Fatal("Error in GetIntegrationPackages: ", err)
	}

	resp.Print()

}

func RunInspectIntegrationPackage(packageName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := InspectIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

func RunGetFlowsOfIntegrationPackage(packageName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := GetFlowsOfIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

func RunDownloadIntegrationPackage(packageName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = DownloadIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in DownloadIntegrationPackage: ", err)
	}

}

func RunInspectFlow(flowName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := InspectFlow(conf, flowName)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print()
}

func RunDownloadFlow(flowId string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = DownloadFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

}

func NewConfiguration() (config.Configuration, error) {
	conf := config.Configuration{}

	// path, _ := os.Getwd()
	// println("os.Getwd(): ", path)

	configFile, err := os.Open("./config.json")
	if err != nil {
		homePath, _ := os.UserHomeDir()
		configFilePath := filepath.Join(homePath, ".cig", "config.json")
		log.Println("Config file path: ", configFilePath)
		configFile, err = os.Open(configFilePath)
		if err != nil {
			err = fmt.Errorf("error opening config.json: %w", err)
			return conf, err
		}

	}

	defer configFile.Close()

	confAll := config.ConfigurationFile{}

	err = json.NewDecoder(configFile).Decode(&confAll)
	if err != nil {
		return conf, err
	}
	for _, v := range confAll.Tenants {
		if confAll.ActiveTenantKey == v.Key {
			conf.Key = v.Key
			conf.ApiURL = v.ApiURL
			conf.Authorization = v.Authorization
			break
		}
	}

	if conf.ApiURL == "" {

		err = errors.New("provide configuration file with tenant ApiURL")
		return conf, err
	}

	return conf, err
}

func RunGetFlowConfigs(flowName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := GetFlowConfigs(conf, flowName)
	if err != nil {
		log.Fatal("Error in GetFlowConfigs: ", err)
	}
	resp.Print()
}

func RunUpdateFlowConfigs(flowName string, configs []model.FlowConfigurationPrinter) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := UpdateFlowConfigsBatch(conf, flowName, configs)
	if err != nil {
		log.Fatal("Error in UpdateFlowConfigs: ", err)
	}

	fmt.Println(resp)
	//resp.Print()

}

func RunCreateFlow(name string, id string, packageid string, fileName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := CreateFlow(conf, name, id, packageid, fileName)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}

	fmt.Println(resp)

}

func RunUpdateFlow(name string, id string, version string, fileName string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := UpdateFlow(conf, name, id, version, fileName)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)

}

func RunDeployFlow(id string, version string) {
	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := DeployFlow(conf, id, version)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)
}
