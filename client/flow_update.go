package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tobiaszgithub/cig/config"
)

//RunUpdateFlow - call the function UpdateFlow
func RunUpdateFlow(conf config.Configuration, name string, id string, version string, fileName string) {

	resp, err := UpdateFlow(conf, name, id, version, fileName)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)

}

//UpdateFlow - update integration flow name and content
func UpdateFlow(conf config.Configuration, name string, id string, version string, fileName string) (string, error) {
	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	var encodedContent string

	if fileName != "" {
		contentData, err := ioutil.ReadFile(fileName)
		if err != nil {
			return "", err
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
			"ArtifactContent": encodedContent,
		}
	} else {
		requestBody = map[string]string{
			"Name": name,
		}
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	updateFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + id + "',Version='" + version + "')"
	log.Println("PUT ", updateFlowURL)

	request, err := http.NewRequest("PUT", updateFlowURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", err
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
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		return "", fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}
	bodyStr := string(body) + "\n"
	return bodyStr, nil
}
