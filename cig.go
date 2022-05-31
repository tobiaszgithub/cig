package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tobiaszgithub/cig/client"
	"github.com/tobiaszgithub/cig/config"
)

func main() {
	println("Cloud Integration CLI")

	conf, err := NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	if conf.ApiURL == "" {
		log.Fatal("Provide configuration file with tenant ApiURL")
	}

	resp, err := client.GetIntegrationPackages(conf)

	if err != nil {
		log.Fatal("Error in GetIntegrationPackages: ", err)
	}

	resp.Print()

}

func NewConfiguration() (config.Configuration, error) {
	configFile, _ := os.Open("config.json")
	defer configFile.Close()

	conf := config.Configuration{}

	confAll := config.ConfigurationFile{}

	err := json.NewDecoder(configFile).Decode(&confAll)

	for _, v := range confAll.Tenants {
		if confAll.ActiveTenantKey == v.Key {
			conf.Key = v.Key
			conf.ApiURL = v.ApiURL
			conf.Authorization = v.Authorization
			break
		}
	}

	return conf, err
}
