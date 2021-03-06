package client

import (
	"bytes"
	"context"
	"encoding/base64"
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
	log.Println(integrationPackagesURL)

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

func InspectIntegrationPackage(conf config.Configuration, packageId string) (*model.IPByIdResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageId + "')"
	log.Println(integrationPackagesURL)

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
	log.Println(flowsOfIntegrationPackagesURL)

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

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

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

func InspectFlow(conf config.Configuration, flowId string) (*model.FlowByIdResponse, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowId + "',Version='active')"
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

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}

func DownloadFlow(conf config.Configuration, flowId string, outputFile string) (string, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowId + "',Version='active')/$value"

	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := GetClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf(string(bodyBytes))
	}

	defer response.Body.Close()

	n, err := saveBodyContent(outputFile, response.Body)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("File created: %s \n", outputFile)
	output += fmt.Sprintf("number of bytes: %d", n)
	return output, nil
}

func saveBodyContent(fileName string, src io.Reader) (writtenBytes int64, err error) {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return 0, err
	// }

	// flowOutPath := filepath.Join(cwd, fileName)

	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	n, err := io.Copy(out, src)
	if err != nil {
		return n, err
	}
	// log.Println(flowOutPath + " created")
	// log.Println("number of bytes: ", n)

	return n, nil
	//body, _ := ioutil.ReadAll(response.Body)
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
	log.Println(createFlowURL)

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

	var decodedRes model.FlowByIdResponse

	//body, _ := ioutil.ReadAll(response.Body)

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
	}

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	//bodyStr := string(body) + "\n"
	return &decodedRes, nil
}

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
		println()
		println(encodedContent)
		println()
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

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	updateFlowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + id + "',Version='" + version + "')"
	log.Println(updateFlowURL)

	request, err := http.NewRequest("PUT", updateFlowURL, bytes.NewBuffer(requestBodyJson))
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

func DeployFlow(conf config.Configuration, id string, version string) (string, error) {

	csrfToken, cookies, err := getCsrfTokenAndCookies(conf)
	if err != nil {
		return "", err
	}

	deployFlowURL := conf.ApiURL + "/DeployIntegrationDesigntimeArtifact?Id='" + id + "'&Version='" + version + "'"
	log.Println(deployFlowURL)

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

func CopyFlow(conf config.Configuration, srcFlowId string, destFlowId string, destFlowName string, destPackageId string) error {

	srcFlow, err := InspectFlow(conf, srcFlowId)
	if err != nil {
		return err
	}

	tmpFileName, err := getTmpFileName()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileName)

	resp, err := DownloadFlow(conf, srcFlowId, tmpFileName)
	if err != nil {
		return err
	}
	log.Print(resp)

	if destFlowName == "" {
		destFlowName = srcFlow.D.Name
	}

	if destPackageId == "" {
		destPackageId = srcFlow.D.PackageID
	}

	createResp, err := CreateFlow(conf, destFlowName, destFlowId, destPackageId, tmpFileName)
	if err != nil {
		return err
	}

	//log.Print(createResp.Print())
	createResp.Print()

	return nil
}

func getTmpFileName() (string, error) {
	tmpfile, err := os.CreateTemp("", "flow*.zip")
	if err != nil {
		return "", err
	}

	err = tmpfile.Close()
	if err != nil {
		return "", err
	}
	tmpFileName := tmpfile.Name()

	err = os.Remove(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return tmpFileName, nil

}
