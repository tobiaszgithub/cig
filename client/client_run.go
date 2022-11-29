package client

import (
	"fmt"
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

//RunUpdateFlowConfigs - call the function UpdateFlowConfigsBatch
func RunUpdateFlowConfigs(conf config.Configuration, flowName string, configs []model.FlowConfigurationPrinter) {

	resp, err := UpdateFlowConfigsBatch(conf, flowName, configs)
	if err != nil {
		log.Fatal("Error in UpdateFlowConfigs: ", err)
	}

	fmt.Println(resp)
	//resp.Print()

}
