package client

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

//RunTransportFlow - call the function TransportFlow
func RunTransportFlow(out io.Writer, conf config.Configuration, srcFlowID string, destFlowID string, destTenantKey string, destFlowName string, destPackageID string) {

	destConf, err := config.NewConfiguration(destTenantKey)
	if err != nil {
		log.Fatal(err)
	}

	err = TransportFlow(out, conf, srcFlowID, destConf, destFlowID, destFlowName, destPackageID)
	if err != nil {
		log.Fatal(err)
	}
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

	outputContent, err := os.OpenFile(tmpFileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Error Openning file:\n", err)
	}
	defer outputContent.Close()

	err = DownloadFlow(out, conf, srcFlowID, version, outputContent)
	if err != nil {
		return err
	}

	if srcFlowID != destFlowID {
		tmpFileName, err = adjustDownloadedFlow(srcFlowID, destFlowID, tmpFileName)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFileName)
	}

	tmpFileContent, err := os.Open(tmpFileName)
	if err != nil {
		return err
	}
	defer tmpFileContent.Close()

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
		err = UpdateFlow(out, destConf, destFlowName, destFlowID, "active", tmpFileName, tmpFileContent)
		fmt.Fprintf(out, "Integration flow updated. Response: %s\n", updateResp)
	} else {
		if destFlowName == "" {
			destFlowName = srcFlow.D.Name
		}

		createResp, err = CreateFlow(destConf, destFlowName, destFlowID, destPackageID, tmpFileContent)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Integration flow created.\n")
		createResp.Print(out)
	}

	if err != nil {
		return err
	}

	return nil
}
