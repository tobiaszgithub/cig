package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/beeekind/go-authhttp"
	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
	"golang.org/x/oauth2/clientcredentials"
)

func GetClient(conf config.Configuration) *http.Client {

	ctx := context.Background()

	var client *http.Client

	if conf.Authorization.Type == "oauth" {
		oauthConf := clientcredentials.Config{
			ClientID:     conf.Authorization.ClientID,
			ClientSecret: conf.Authorization.ClientSecret,
			Scopes:       []string{},
			TokenURL:     conf.Authorization.TokenURL,
		}

		client = oauthConf.Client(ctx)
	}

	if conf.Authorization.Type == "basic" {
		client = authhttp.NewHTTPClient(authhttp.WithBasicAuth(conf.Authorization.Username,
			conf.Authorization.Password))
	}

	return client

}

func GetIntegrationPackages(conf config.Configuration) (*model.IPResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages"

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	rawRes, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer rawRes.Body.Close()

	var decodedRes model.IPResponse

	if err := json.NewDecoder(rawRes.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err

}

func InspectIntegrationPackage(conf config.Configuration, packageName string) (*model.IPByIdResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageName + "')"

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	rawRes, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer rawRes.Body.Close()

	var decodedRes model.IPByIdResponse

	if err := json.NewDecoder(rawRes.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}

func GetFlowsOfIntegrationPackage(conf config.Configuration, packageName string) (*model.FlowsOfIPResponse, error) {
	flowsOfIntegrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageName + "')/IntegrationDesigntimeArtifacts"

	request, err := http.NewRequest("GET", flowsOfIntegrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var decodedRes model.FlowsOfIPResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err

}

func DownloadIntegrationPackage(conf config.Configuration, packageName string) error {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageName + "')/$value"

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(bodyBytes))
	}

	defer response.Body.Close()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	packageOutPath := filepath.Join(cwd, packageName+".zip")
	//log.Println(packageOutPath)

	out, err := os.Create(packageOutPath)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	log.Println(packageOutPath + " created")
	log.Println("number of bytes: ", n)

	return nil
}

func InspectFlow(conf config.Configuration, fileName string) (*model.FlowByIdResponse, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + fileName + "',Version='active')"
	log.Println(flowURL)
	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var decodedRes model.FlowByIdResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}

func DownloadFlow(conf config.Configuration, flowName string) error {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='active')/$value"

	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(bodyBytes))
	}

	defer response.Body.Close()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	flowOutPath := filepath.Join(cwd, flowName+".zip")
	//log.Println(packageOutPath)

	out, err := os.Create(flowOutPath)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	log.Println(flowOutPath + " created")
	log.Println("number of bytes: ", n)

	return nil
}

func GetConfigsFlow(conf config.Configuration, fileName string) (*model.FlowConfigurations, error) {
	configsFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + fileName + "',Version='active')/Configurations"
	log.Println(configsFlowURL)
	request, err := http.NewRequest("GET", configsFlowURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var decodedRes model.FlowConfigurations

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}
