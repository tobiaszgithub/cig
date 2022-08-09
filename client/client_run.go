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

func RunInspectFlow(conf config.Configuration, flowId string) {
	// conf, err := config.NewDefaultConfiguration()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := InspectFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print(os.Stdout)
}

func RunDownloadFlow(conf config.Configuration, flowId string, outputFile string) {
	if outputFile == "" {
		outputFile = flowId + ".zip"
	}

	resp, err := DownloadFlow(conf, flowId, outputFile)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

	fmt.Println(resp)

}

func RunGetFlowConfigs(conf config.Configuration, flowName string, fileName string) {

	resp, err := GetFlowConfigs(conf, flowName)
	if err != nil {
		log.Fatal("Error in GetFlowConfigs: ", err)
	}

	resp.Print(os.Stdout)

	var outputFile *os.File

	if fileName != "" {
		log.Println("File name: ", fileName)
		outputFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal("Error creating file: ", err)
		}
		defer outputFile.Close()
	}

	if outputFile != nil {
		resp.Print(outputFile)
	}
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

func RunDeployFlow(conf config.Configuration, id string, version string) {

	resp, err := DeployFlow(conf, id, version)
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
