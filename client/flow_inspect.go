package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

func RunInspectFlow(conf config.Configuration, flowId string, version string) {

	resp, err := InspectFlow(conf, flowId, version)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print(os.Stdout)
}

func InspectFlow(conf config.Configuration, flowId string, version string) (*model.FlowByIdResponse, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowId + "',Version='active')"
	log.Println("GET ", flowURL)
	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	defer response.Body.Close()

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
		//return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, body)
	}

	var decodedRes model.FlowByIdResponse
	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}
