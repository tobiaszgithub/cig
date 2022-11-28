package client

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/tobiaszgithub/cig/config"
)

//RunCopyFlow - call the function CopyFlow
func RunCopyFlow(conf config.Configuration, srcFlowID string, destFlowID string, destFlowName string, destPackageID string) {
	err := CopyFlow(conf, srcFlowID, destFlowID, destFlowName, destPackageID)
	if err != nil {
		log.Fatal(err)
	}
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

	outputContent, err := os.OpenFile(tmpFileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Error Openning file:\n", err)
	}
	defer outputContent.Close()

	var out bytes.Buffer
	err = DownloadFlow(&out, conf, srcFlowID, version, tmpFileName, outputContent)
	if err != nil {
		return err
	}
	fmt.Println(out.String())

	if destFlowName == "" {
		destFlowName = srcFlow.D.Name
	}

	if destPackageID == "" {
		destPackageID = srcFlow.D.PackageID
	}

	tmpFileContent, err := os.Open(tmpFileName)
	if err != nil {
		return err
	}
	defer tmpFileContent.Close()

	createResp, err := CreateFlow(conf, destFlowName, destFlowID, destPackageID, tmpFileContent)
	if err != nil {
		return err
	}

	//log.Print(createResp.Print())
	createResp.Print(os.Stdout)

	return nil
}
