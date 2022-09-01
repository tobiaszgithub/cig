package client

import (
	"fmt"
	"io"
	"log"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

//RunGetIntegrationPackages - call the GetIntegrationPackages
func RunGetIntegrationPackages(conf config.Configuration) {

	// conf, err := config.NewConfiguration(tenantKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := GetIntegrationPackages(conf)

	if err != nil {
		log.Fatal("Error in GetIntegrationPackages: ", err)
	}

	resp.Print()

}

//RunInspectIntegrationPackage - call the InspectIntegrationPackage
func RunInspectIntegrationPackage(conf config.Configuration, packageID string) {
	// conf, err := config.NewDefaultConfiguration()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := InspectIntegrationPackage(conf, packageID)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

//RunGetFlowsOfIntegrationPackage - call the GetFlowsOfIntegrationPackage
func RunGetFlowsOfIntegrationPackage(conf config.Configuration, packageName string) {
	// conf, err := config.NewDefaultConfiguration()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := GetFlowsOfIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

//RunDownloadIntegrationPackage - call the function DownloadIntegrationPackage
func RunDownloadIntegrationPackage(conf config.Configuration, packageName string) {
	// conf, err := config.NewDefaultConfiguration()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := DownloadIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in DownloadIntegrationPackage: ", err)
	}

}

//RunDownloadFlow - call the function DownloadFlow
func RunDownloadFlow(conf config.Configuration, flowID string, version string, outputFile string) {
	if outputFile == "" {
		outputFile = flowID + ".zip"
	}

	resp, err := DownloadFlow(conf, flowID, version, outputFile)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

	fmt.Println(resp)

}

//RunUpdateFlowConfigs - call the function UpdateFlowConfigsBatch
func RunUpdateFlowConfigs(conf config.Configuration, flowName string, configs []model.FlowConfigurationPrinter) {

	resp, err := UpdateFlowConfigsBatch(conf, flowName, configs)
	if err != nil {
		log.Fatal("Error in UpdateFlowConfigs: ", err)
	}

	fmt.Println(resp)
	//resp.Print()

}

//RunCopyFlow - call the function CopyFlow
func RunCopyFlow(conf config.Configuration, srcFlowID string, destFlowID string, destFlowName string, destPackageID string) {
	err := CopyFlow(conf, srcFlowID, destFlowID, destFlowName, destPackageID)
	if err != nil {
		log.Fatal(err)
	}
}

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
