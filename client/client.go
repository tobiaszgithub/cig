package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func GetFlowConfigs(conf config.Configuration, flowName string) (*model.FlowConfigurations, error) {
	configsFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='active')/Configurations"
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

func UpdateFlowConfigs(conf config.Configuration, flowName string, configs []model.FlowConfigurationPrinter) (string, error) {

	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	var bodyStr string
	for _, c := range configs {
		//key, value, err := parseParam(param)
		if err != nil {
			return "", err
		}

		requestBody := map[string]string{"ParameterValue": c.ParameterValue, "DataType": c.DataType}
		requestBodyJson, err := json.Marshal(requestBody)
		if err != nil {
			return "", err
		}
		updateFlowConfigsURL := conf.ApiURL +
			"/IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='active')/$links/Configurations('" + c.ParameterKey + "')"
		log.Println(updateFlowConfigsURL)

		request, err := http.NewRequest("PUT", updateFlowConfigsURL, bytes.NewBuffer(requestBodyJson))
		if err != nil {
			return "", err
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
			return "", err
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		//	fmt.Println("response Body:", string(body))

		statusOk := response.StatusCode >= 200 && response.StatusCode < 300
		if !statusOk {
			return "", fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
		}
		bodyStr = bodyStr + string(body) + "\n"

	}
	return bodyStr, nil
}

func UpdateFlowConfigsBatch(conf config.Configuration, flowName string, configs []model.FlowConfigurationPrinter) (string, error) {

	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	var bodyStr string

	batchURL := conf.ApiURL + "/$batch"
	method := "POST"

	begining :=
		"--batch_request\r\n" +
			"Content-Type: multipart/mixed; boundary=changeset_abc\r\n" +
			"\r\n"

	end :=
		"--changeset_abc--\r\n" +
			"\r\n" +
			"--batch_request--"

	var payload string
	var batch string
	for _, c := range configs {

		updateFlowConfigsURL := //conf.ApiURL +
			"IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='active')/$links/Configurations('" + c.ParameterKey + "')"
		log.Println(updateFlowConfigsURL)

		singleRequest :=
			"--changeset_abc\r\n" +
				"Content-Type: application/http\r\n" +
				"Content-Transfer-Encoding:binary\r\n" +
				"\r\n" +
				"PUT " +

				//"IntegrationDesigntimeArtifacts(Id='PurchaseOrder',Version='active')/$links/Configurations('bodySize')" +
				updateFlowConfigsURL +
				" HTTP/1.1\r\n" +
				"Accept: application/json\r\n" +
				"Content-Type: application/json\r\n" +
				"\r\n" +
				"{\r\n" +
				"\"ParameterValue\": \"" + c.ParameterValue + "\",\r\n" +
				"\"DataType\": \"" + c.DataType + "\"\r\n" +
				"}\r\n" +
				"\r\n"

		batch = batch + singleRequest

	}
	payload = begining + batch + end
	println(payload)
	request, err := http.NewRequest(method, batchURL, bytes.NewBufferString(payload))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "multipart/mixed; boundary=batch_request")
	request.Header.Set("X-CSRF-Token", csrfToken)
	//request.Header.Set("Content-Length", "431")
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
	//	fmt.Println("response Body:", string(body))

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		return "", fmt.Errorf("response Status: %s\n, response body:\n %s", response.Status, string(body))
	}
	bodyStr = bodyStr + string(body) + "\n"

	return bodyStr, nil
}

func getCsrfTokenAndCookies(conf config.Configuration) (string, []*http.Cookie, error) {
	csrfTokenURL := conf.ApiURL + "/"
	tokenRequest, err := http.NewRequest("GET", csrfTokenURL, nil)
	if err != nil {
		return "", nil, err
	}
	tokenRequest.Header.Set("X-CSRF-Token", "Fetch")
	tokenHttpClient := GetClient(conf)

	tokenResponse, err := tokenHttpClient.Do(tokenRequest)
	if err != nil {
		return "", nil, err
	}
	defer tokenResponse.Body.Close()
	csrfToken := tokenResponse.Header.Get("X-CSRF-Token")
	cookies := tokenResponse.Cookies()

	return csrfToken, cookies, nil
}

func CreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string) (string, error) {
	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	requestBody := map[string]string{"Name": name,
		"Id":        id,
		"PackageId": packageid}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	createFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts"
	log.Println(createFlowURL)

	request, err := http.NewRequest("POST", createFlowURL, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return "", err
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
