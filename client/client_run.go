package client

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tobiaszgithub/cig/config"
	"github.com/tobiaszgithub/cig/model"
)

func RunGetIntegrationPackages() {

	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := GetIntegrationPackages(conf)

	if err != nil {
		log.Fatal("Error in GetIntegrationPackages: ", err)
	}

	resp.Print()

}

func RunInspectIntegrationPackage(packageId string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := InspectIntegrationPackage(conf, packageId)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

func RunGetFlowsOfIntegrationPackage(packageName string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := GetFlowsOfIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in InspectIntegrationPackage: ", err)
	}
	resp.Print()

}

func RunDownloadIntegrationPackage(packageName string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = DownloadIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in DownloadIntegrationPackage: ", err)
	}

}

func RunInspectFlow(flowId string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := InspectFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print()
}

func RunDownloadFlow(flowId string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = DownloadFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

}

func RunGetFlowConfigs(flowName string, outputFile io.WriteSeeker) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := GetFlowConfigs(conf, flowName)
	if err != nil {
		log.Fatal("Error in GetFlowConfigs: ", err)
	}

	resp.Print(os.Stdout)

	if outputFile != nil {
		resp.Print(outputFile)
	}
}

func RunUpdateFlowConfigs(flowName string, configs []model.FlowConfigurationPrinter) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := UpdateFlowConfigsBatch(conf, flowName, configs)
	if err != nil {
		log.Fatal("Error in UpdateFlowConfigs: ", err)
	}

	fmt.Println(resp)
	//resp.Print()

}

func RunCreateFlow(name string, id string, packageid string, fileName string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := CreateFlow(conf, name, id, packageid, fileName)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}

	fmt.Println(resp)

}

func RunUpdateFlow(name string, id string, version string, fileName string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := UpdateFlow(conf, name, id, version, fileName)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)

}

func RunDeployFlow(id string, version string) {
	conf, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := DeployFlow(conf, id, version)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)
}
