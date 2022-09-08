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

//RunInspectFlow - call the function InspectFlow
func RunInspectFlow(conf config.Configuration, flowID string, version string) {

	resp, err := InspectFlow(conf, flowID, version)
	if err != nil {
		log.Fatal("Error in InspectFlow:\n", err)
	}
	resp.Print(os.Stdout)
}

//InspectFlow - inspect flow
func InspectFlow(conf config.Configuration, flowID string, version string) (*model.FlowByIdResponse, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowID + "',Version='active')"
	log.Println("GET ", flowURL)
	request, err := http.NewRequest("GET", flowURL, nil)
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
