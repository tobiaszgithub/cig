package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tobiaszgithub/cig/config"
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
		log.Fatal("Error in InspectIntegrationPackage: ", err)
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
