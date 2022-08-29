package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

func RunCreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string) {
	resp, err := CreateFlow(conf, name, id, packageid, fileName)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}
	//fmt.Println(resp)
	resp.Print(os.Stdout)

}

func CreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string) (*model.FlowByIdResponse, error) {
	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return nil, err
	}

	var encodedContent string

	if fileName != "" {
		contentData, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		encodedContent = base64.StdEncoding.EncodeToString(contentData)
		// println()
		// println(encodedContent)
		// println()
	}

	var requestBody map[string]string

	if encodedContent != "" {
		requestBody = map[string]string{
			"Name":            name,
			"Id":              id,
			"PackageId":       packageid,
			"ArtifactContent": encodedContent,
		}
	} else {
		requestBody = map[string]string{
			"Name":      name,
			"Id":        id,
			"PackageId": packageid,
		}
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	createFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts"
	log.Println("POST ", createFlowURL)

	request, err := http.NewRequest("POST", createFlowURL, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-CSRF-Token", csrfToken)
	for i := range cookies {
		request.AddCookie(cookies[i])
	}

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//body, _ := ioutil.ReadAll(response.Body)

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

	var decodedRes model.FlowByIdResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	//bodyStr := string(body) + "\n"
	return &decodedRes, nil
}
