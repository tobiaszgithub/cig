package client

import (
	"bytes"
	"encoding/base64"
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

//RunCreateFlow - call the function CreateFlow
func RunCreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string) {

	var fileContent *os.File
	if fileName != "" {
		fileContent, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Error Openning file: ", err)
		}
		defer fileContent.Close()
	} else {
		fileContent = nil
	}

	resp, err := CreateFlow(conf, name, id, packageid, fileName, fileContent)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}
	//fmt.Println(resp)
	resp.Print(os.Stdout)

}

//CreateFlow - create integration flow, it is possible to create empty integration flow or with content
func CreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string, flowContent io.Reader) (*model.FlowByIdResponse, error) {
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

	if flowContent != nil {
		contentData, err := io.ReadAll(flowContent)
		if err != nil {
			return nil, err
		}
		encodedContent = base64.StdEncoding.EncodeToString(contentData)

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

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	createFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts"
	log.Println("POST ", createFlowURL)

	request, err := http.NewRequest("POST", createFlowURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-CSRF-Token", csrfToken)
	for i := range cookies {
		request.AddCookie(cookies[i])
	}

	httpClient := getClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer response.Body.Close()

	//body, _ := ioutil.ReadAll(response.Body)

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
		//return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

	var decodedRes model.FlowByIdResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	//bodyStr := string(body) + "\n"
	return &decodedRes, nil
}
