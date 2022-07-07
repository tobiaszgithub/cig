package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ConfigurationFile struct {
	ActiveTenantKey string
	Tenants         []struct {
		Key           string
		ApiURL        string
		Authorization struct {
			Type         string
			Username     string
			Password     string
			ClientID     string
			ClientSecret string
			TokenURL     string
		}
	}
}

type Configuration struct {
	Key           string
	ApiURL        string
	Authorization struct {
		Type         string
		Username     string
		Password     string
		ClientID     string
		ClientSecret string
		TokenURL     string
	}
}

func NewConfiguration() (Configuration, error) {
	conf := Configuration{}

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

	confAll := ConfigurationFile{}

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
