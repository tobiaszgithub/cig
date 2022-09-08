package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

//RunGetFlowConfigs - call the function GetFlowConfigs
func RunGetFlowConfigs(out io.Writer, conf config.Configuration, flowName string, fileName string, version string) {

	resp, err := GetFlowConfigs(conf, flowName, version)
	if err != nil {
		log.Fatal("Error in GetFlowConfigs:\n", err)
	}

	resp.Print(out)

	var outputFile *os.File

	if fileName != "" {
		log.Println("File name: ", fileName)
		outputFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal("Error creating file:\n", err)
		}
		defer outputFile.Close()
	}

	if outputFile != nil {
		resp.Print(outputFile)
	}
}

//GetFlowConfigs - get integration flows configuration
func GetFlowConfigs(conf config.Configuration, flowName string, version string) (*model.FlowConfigurations, error) {
	configsFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='" + version + "')/Configurations"
	log.Println("GET ", configsFlowURL)
	request, err := http.NewRequest("GET", configsFlowURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer response.Body.Close()
	log.Printf("response statusCode: %d\n", response.StatusCode)
	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if response.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, body)
	}

	var decodedRes model.FlowConfigurations

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}
