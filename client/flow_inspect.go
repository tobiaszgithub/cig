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

func RunInspectFlow(conf config.Configuration, flowId string) {

	resp, err := InspectFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print(os.Stdout)
}

func InspectFlow(conf config.Configuration, flowId string) (*model.FlowByIdResponse, error) {
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

	var decodedRes model.FlowByIdResponse

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

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}

func DeployFlow(conf config.Configuration, id string, version string) (string, error) {

	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	deployFlowURL := conf.ApiURL + "/DeployIntegrationDesigntimeArtifact?Id='" + id + "'&Version='" + version + "'"
	log.Println("POST ", deployFlowURL)

	request, err := http.NewRequest("POST", deployFlowURL, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", "application/json")

	request.Header.Set("X-CSRF-Token", csrfToken)
	for i := range cookies {
		request.AddCookie(cookies[i])
	}

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		return "", fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}
	bodyStr := "Task ID:\n" + string(body) + "\n"
	return bodyStr, nil
}
