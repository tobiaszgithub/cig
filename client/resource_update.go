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
	"path/filepath"

	"github.com/tobiaszgithub/cig/config"
)

//RunResourceUpdate - call ResourceUpdate
func RunResourceUpdate(out io.Writer, conf config.Configuration, flowID string, flowVersion string, resourceName string, resourceType string, resourceFileName string) {
	err := ResourceUpdate(out, conf, flowID, flowVersion, resourceName, resourceType, resourceFileName)
	if err != nil {
		log.Fatal(err)
	}
}

//ResourceUpdate - update resource of the integration flow
func ResourceUpdate(out io.Writer, conf config.Configuration, flowID string, flowVersion string, resourceName string, resourceType string, resourceFileName string) error {
	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return err
	}

	if resourceName == "" {
		resourceName = filepath.Base(resourceFileName)
	}

	var encodedContent string
	contentData, err := ioutil.ReadFile(resourceFileName)
	if err != nil {
		return err
	}
	encodedContent = base64.StdEncoding.EncodeToString(contentData)

	requestBody := map[string]string{
		"ResourceContent": encodedContent,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	updateResourceURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowID + "',Version='" + flowVersion + "')" +
		"/$links/Resources(Name='" + resourceName + "',ResourceType='" + resourceType + "')"
	log.Println("PUT ", updateResourceURL)

	request, err := http.NewRequest("PUT", updateResourceURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return err
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
		return err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		return fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

	return nil

}
