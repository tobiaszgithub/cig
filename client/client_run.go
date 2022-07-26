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

	conf, err := config.NewDefaultConfiguration()
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
	conf, err := config.NewDefaultConfiguration()
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
	conf, err := config.NewDefaultConfiguration()
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
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = DownloadIntegrationPackage(conf, packageName)
	if err != nil {
		log.Fatal("Error in DownloadIntegrationPackage: ", err)
	}

}

func RunInspectFlow(flowId string) {
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := InspectFlow(conf, flowId)
	if err != nil {
		log.Fatal("Error in InspectFlow: ", err)
	}
	resp.Print(os.Stdout)
}

func RunDownloadFlow(flowId string, outputFile string) {
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	if outputFile == "" {
		outputFile = flowId + ".zip"
	}

	resp, err := DownloadFlow(conf, flowId, outputFile)
	if err != nil {
		log.Fatal("Error in DownloadFlow: ", err)
	}

	fmt.Println(resp)

}

func RunGetFlowConfigs(flowName string, fileName string) {
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

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

func RunUpdateFlowConfigs(flowName string, configs []model.FlowConfigurationPrinter) {
	conf, err := config.NewDefaultConfiguration()
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
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := CreateFlow(conf, name, id, packageid, fileName)
	if err != nil {
		log.Fatal("Error in CreateFlow:\n", err)
	}
	//fmt.Println(resp)
	resp.Print(os.Stdout)

}

func RunUpdateFlow(name string, id string, version string, fileName string) {
	conf, err := config.NewDefaultConfiguration()
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
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := DeployFlow(conf, id, version)
	if err != nil {
		log.Fatal("Error in UpdateFlow:\n", err)
	}

	fmt.Println(resp)
}

func RunCopyFlow(srcFlowId string, destFlowId string, destFlowName string, destPackageId string) {
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	err = CopyFlow(conf, srcFlowId, destFlowId, destFlowName, destPackageId)
	if err != nil {
		log.Fatal(err)
	}
}

func RunTransportFlow(out io.Writer, srcFlowId string, destFlowId string, destTenantKey string, destFlowName string, destPackageId string) {
	conf, err := config.NewDefaultConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	destConf, err := config.NewConfiguration(destTenantKey)
	if err != nil {
		log.Fatal(err)
	}

	err = TransportFlow(out, conf, srcFlowId, destConf, destFlowId, destFlowName, destPackageId)
	if err != nil {
		log.Fatal(err)
	}
}
