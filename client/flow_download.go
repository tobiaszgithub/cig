package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tobiaszgithub/cig/config"
)

//RunDownloadFlow - call the function DownloadFlow
func RunDownloadFlow(out io.Writer, conf config.Configuration, flowID string, version string, outputFile string) {
	if outputFile == "" {
		outputFile = flowID + ".zip"
	}

	outputContent, err := os.OpenFile(outputFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Error Openning file:\n", err)
	}

	err = DownloadFlow(out, conf, flowID, version, outputContent)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

}

//DownloadFlow is the function to download integration flow content
func DownloadFlow(out io.Writer, conf config.Configuration, flowID string, version string, outputContent io.Writer) error {
	flowURL := conf.ApiURL + "/IntegrationDesigntimeArtifacts(Id='" + flowID + "',Version='active')/$value"
	log.Println("GET ", flowURL)
	request, err := http.NewRequest("GET", flowURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")

	httpClient := getClient(conf)

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer response.Body.Close()

	statusOk := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOk {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if response.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, body)
	}

	n, err := saveBodyContent(outputContent, response.Body)
	if err != nil {
		return err
	}

	output := fmt.Sprintf("Content downloaded.\n")
	output += fmt.Sprintf("number of bytes: %d", n)
	fmt.Fprintf(out, "%s", output)
	return nil
}

func saveBodyContent(outputContent io.Writer, src io.Reader) (writtenBytes int64, err error) {

	// out, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	// if err != nil {
	// 	return 0, err
	// }
	// defer out.Close()

	n, err := io.Copy(outputContent, src)
	if err != nil {
		return n, err
	}
	// log.Println(flowOutPath + " created")
	// log.Println("number of bytes: ", n)

	return n, nil
	//body, _ := ioutil.ReadAll(response.Body)
}
