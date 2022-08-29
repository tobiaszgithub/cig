package client

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

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

func RunInspectIntegrationPackage(conf config.Configuration, packageId string) {
	// conf, err := config.NewDefaultConfiguration()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := InspectIntegrationPackage(conf, packageId)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

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

func RunDownloadFlow(conf config.Configuration, flowId string, version string, outputFile string) {
	if outputFile == "" {
		outputFile = flowId + ".zip"
	}

	resp, err := DownloadFlow(conf, flowId, version, outputFile)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

	fmt.Println(resp)

}

func RunUpdateFlowConfigs(conf config.Configuration, flowName string, configs []model.FlowConfigurationPrinter) {

	resp, err := UpdateFlowConfigsBatch(conf, flowName, configs)
	if err != nil {
		log.Fatal("Error in UpdateFlowConfigs: ", err)
	}

	fmt.Println(resp)
	//resp.Print()

}

func RunCreateFlow(conf config.Configuration, name string, id string, packageid string, fileName string) {
	resp, err := CreateFlow(conf, name, id, packageid, fileName)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}
	//fmt.Println(resp)
	resp.Print(os.Stdout)

}

func RunUpdateFlow(conf config.Configuration, name string, id string, version string, fileName string) {

	resp, err := UpdateFlow(conf, name, id, version, fileName)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)

}

func RunCopyFlow(conf config.Configuration, srcFlowId string, destFlowId string, destFlowName string, destPackageId string) {
	err := CopyFlow(conf, srcFlowId, destFlowId, destFlowName, destPackageId)
	if err != nil {
		log.Fatal(err)
	}
}

func RunTransportFlow(out io.Writer, conf config.Configuration, srcFlowId string, destFlowId string, destTenantKey string, destFlowName string, destPackageId string) {

	destConf, err := config.NewConfiguration(destTenantKey)
	if err != nil {
		log.Fatal(err)
	}

	err = TransportFlow(out, conf, srcFlowId, destConf, destFlowId, destFlowName, destPackageId)
	if err != nil {
		log.Fatal(err)
	}
}
