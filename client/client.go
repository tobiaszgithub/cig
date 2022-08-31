package client

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/beeekind/go-authhttp"
	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	//ErrConnection - connection error
	ErrConnection = errors.New("connection error")
	//ErrNotFound - not found
	ErrNotFound = errors.New("not found")
	//ErrInvalidResponse - invalid server response
	ErrInvalidResponse = errors.New("invalid server response")
	//ErrInvalid - invalid data
	ErrInvalid = errors.New("invalid data")
	//ErrNotNumber - not a number
	ErrNotNumber = errors.New("not a number")
)

func getClient(conf config.Configuration) *http.Client {

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

//GetIntegrationPackages is the function to get list of the integration packages
func GetIntegrationPackages(conf config.Configuration) (*model.IPResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages"
	log.Println("GET ", integrationPackagesURL)

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

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

	var decodedRes model.IPResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err

}

//InspectIntegrationPackage is the function to get details of the integration package
func InspectIntegrationPackage(conf config.Configuration, packageID string) (*model.IPByIdResponse, error) {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageID + "')"
	log.Println("GET ", integrationPackagesURL)

	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

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

	var decodedRes model.IPByIdResponse

	if err := json.NewDecoder(response.Body).Decode(&decodedRes); err != nil {
		return nil, err
	}

	return &decodedRes, err
}

//GetFlowsOfIntegrationPackage is the function to get list of integration flow of the integration package
func GetFlowsOfIntegrationPackage(conf config.Configuration, packageName string) (*model.FlowsOfIPResponse, error) {
	flowsOfIntegrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageName + "')/IntegrationDesigntimeArtifacts"
	log.Println("GET ", flowsOfIntegrationPackagesURL)

	request, err := http.NewRequest("GET", flowsOfIntegrationPackagesURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

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

//DownloadIntegrationPackage is the function to download all content of the integration package. The objects
//in the integration package have to be in final state
func DownloadIntegrationPackage(conf config.Configuration, packageName string) error {
	integrationPackagesURL := conf.ApiURL + "/IntegrationPackages('" + packageName + "')/$value"
	log.Println("GET ", integrationPackagesURL)
	request, err := http.NewRequest("GET", integrationPackagesURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

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

	n, err := saveBodyContent(packageOutPath, response.Body)
	if err != nil {
		return err
	}

	// out, err := os.Create(packageOutPath)
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

	// n, err := io.Copy(out, response.Body)
	// if err != nil {
	// 	return err
	// }
	log.Println(packageOutPath + " created")
	log.Println("number of bytes: ", n)

	return nil
}

//DownloadFlow is the function to download integration flow content
func DownloadFlow(conf config.Configuration, flowID string, version string, outputFile string) (string, error) {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowID + "',Version='active')/$value"
	log.Println("GET ", flowURL)
	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

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

//UpdateFlowConfigs is the function to update flow's configuration
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
		requestBodyJSON, err := json.Marshal(requestBody)
		if err != nil {
			return "", err
		}
		updateFlowConfigsURL := conf.ApiURL +
			"/IntegrationDesigntimeArtifacts(Id='" + flowName + "',Version='active')/$links/Configurations('" + c.ParameterKey + "')"
		log.Println("PUT ", updateFlowConfigsURL)

		request, err := http.NewRequest("PUT", updateFlowConfigsURL, bytes.NewBuffer(requestBodyJSON))
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
		//	fmt.Println("response Body:", string(body))

		statusOk := response.StatusCode >= 200 && response.StatusCode < 300
		if !statusOk {
			return "", fmt.Errorf("response Status: %s, response body: %s", response.Status, string(body))
		}
		bodyStr = bodyStr + string(body) + "\n"

	}
	return bodyStr, nil
}

//UpdateFlowConfigsBatch is the function to update flow's configuration
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
		log.Println("POST ", updateFlowConfigsURL)

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

	httpClient := getClient(conf)

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
	tokenHTTPClient := getClient(conf)

	tokenResponse, err := tokenHTTPClient.Do(tokenRequest)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer tokenResponse.Body.Close()
	csrfToken := tokenResponse.Header.Get("X-CSRF-Token")
	cookies := tokenResponse.Cookies()

	return csrfToken, cookies, nil
}

//CreateFlow is the function to create integration flow. Integration flow content can be empty
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

//UpdateFlow is the function to update flow name and content
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

//CopyFlow is the function to copy flows in the same system
func CopyFlow(conf config.Configuration, srcFlowID string, destFlowID string, destFlowName string, destPackageID string) error {
	version := "active"
	srcFlow, err := InspectFlow(conf, srcFlowID, version)
	if err != nil {
		return err
	}

	tmpFileName, err := getTmpFileName()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileName)

	resp, err := DownloadFlow(conf, srcFlowID, version, tmpFileName)
	if err != nil {
		return err
	}
	log.Print(resp)

	if destFlowName == "" {
		destFlowName = srcFlow.D.Name
	}

	if destPackageID == "" {
		destPackageID = srcFlow.D.PackageID
	}

	createResp, err := CreateFlow(conf, destFlowName, destFlowID, destPackageID, tmpFileName)
	if err != nil {
		return err
	}

	//log.Print(createResp.Print())
	createResp.Print(os.Stdout)

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

//TransportFlow is the function for Transporting flow from one system to another
func TransportFlow(out io.Writer, conf config.Configuration, srcFlowID string, destConf config.Configuration, destFlowID string, destFlowName string, destPackageID string) error {
	version := "active"
	srcFlow, err := InspectFlow(conf, srcFlowID, version)
	if err != nil {
		return err
	}

	tmpFileName, err := getTmpFileName()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileName)

	resp, err := DownloadFlow(conf, srcFlowID, version, tmpFileName)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", resp)

	if srcFlowID != destFlowID {
		tmpFileName, err = adjustDownloadedFlow(srcFlowID, destFlowID, tmpFileName)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFileName)
	}

	if destPackageID == "" {
		destPackageID = srcFlow.D.PackageID
	}

	destFlow, _ := InspectFlow(destConf, destFlowID, version)

	var createResp *model.FlowByIdResponse
	var updateResp string
	if destFlow != nil && destFlow.D.ID != "" {
		if destFlowName == "" {
			destFlowName = destFlow.D.Name
		}
		updateResp, err = UpdateFlow(destConf, destFlowName, destFlowID, "active", tmpFileName)
		fmt.Fprintf(out, "Integration flow updated. Response: %s\n", updateResp)
	} else {
		if destFlowName == "" {
			destFlowName = srcFlow.D.Name
		}

		createResp, err = CreateFlow(destConf, destFlowName, destFlowID, destPackageID, tmpFileName)
		fmt.Fprintf(out, "Integration flow created.\n")
		createResp.Print(out)
	}

	if err != nil {
		return err
	}

	return nil
}

func adjustDownloadedFlow(srcFlowID, destFlowID string, zipFile string) (newZipFile string, err error) {

	fileDestinationFolder := zipFile[:len(zipFile)-4]

	err = unzipFile(zipFile, fileDestinationFolder)
	if err != nil {
		return "", fmt.Errorf("error during uzipping file: %w", err)
	}
	log.Printf("file: %s has been unzipped to directory %s\n", zipFile, fileDestinationFolder)

	manifestFile := filepath.Join(fileDestinationFolder, "META-INF/MANIFEST.MF")
	oldValue := "SymbolicName: " + srcFlowID
	newValue := "SymbolicName: " + destFlowID
	err = replaceFileContent(manifestFile, oldValue, newValue)
	if err != nil {
		return "", fmt.Errorf("error updating META-INF/MANIFEST.MF file: %w", err)
	}

	projectFile := filepath.Join(fileDestinationFolder, ".project")
	oldValue = "<name>" + srcFlowID + "</name>"
	newValue = "<name>" + destFlowID + "</name>"
	err = replaceFileContent(projectFile, oldValue, newValue)
	if err != nil {
		return "", fmt.Errorf("error updating .project file: %w", err)
	}

	newZipFile = fileDestinationFolder + "Copy.zip"
	if err := zipSource(fileDestinationFolder+string(filepath.Separator), newZipFile); err != nil {
		return "", fmt.Errorf("error creating new zip file: %w", err)
	}
	log.Printf("new zip file has been created: %s", newZipFile)

	return newZipFile, nil
}

func replaceFileContent(filename string, old string, new string) error {

	oldContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	newContents := strings.Replace(string(oldContents), old, new, -1)
	err = ioutil.WriteFile(filename, []byte(newContents), fs.ModePerm)
	if err != nil {
		return fmt.Errorf("error updating file: %w", err)
	}
	log.Printf("File: %s has been updated.\n", filename)

	return nil
}

func unzipFile(sourceFile, targetDirectory string) (err error) {

	openedFile, err := zip.OpenReader(sourceFile)
	if err != nil {
		return err
	}
	defer openedFile.Close()

	for _, file := range openedFile.File {
		if strings.Contains(file.Name, "..") {
			continue
		}
		filePath := filepath.Join(targetDirectory, file.Name)
		//log.Println("unzipping file", filePath)
		if file.FileInfo().IsDir() {
			// create the directory
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		} else {
			destinationFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
			if err != nil {
				panic(err)
			}
			//Opening the file and copy it's contents
			fileInArchive, err := file.Open()

			if err != nil {
				return err
			}

			if _, err := io.Copy(destinationFile, fileInArchive); err != nil {
				panic(err)
			}

			destinationFile.Close()
			fileInArchive.Close()
		}

	}

	return nil
}

func zipSource(sourceDirectory, targetFile string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		//log.Println(path)
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(sourceDirectory), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
