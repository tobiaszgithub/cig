package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tobiaszgithub/cig/config"
)

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
