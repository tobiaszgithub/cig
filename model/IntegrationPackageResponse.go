package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lensesio/tableprinter"
)

type IntegrationPackage struct {
	ID                string `json:"Id"`
	Name              string `json:"Name"`
	Description       string `json:"Description"`
	ShortText         string `json:"ShortText"`
	Version           string `json:"Version"`
	Vendor            string `json:"Vendor"`
	PartnerContent    bool   `json:"PartnerContent"`
	UpdateAvailable   bool   `json:"UpdateAvailable"`
	Mode              string `json:"Mode"`
	SupportedPlatform string `json:"SupportedPlatform"`
	ModifiedBy        string `json:"ModifiedBy"`
	CreationDate      string `json:"CreationDate"`
	ModifiedDate      string `json:"ModifiedDate"`
	CreatedBy         string `json:"CreatedBy"`
	Products          string `json:"Products"`
	Keywords          string `json:"Keywords"`
	Countries         string `json:"Countries"`
	Industries        string `json:"Industries"`
	LineOfBusiness    string `json:"LineOfBusiness"`
}

type IPResponse struct {
	D struct {
		Results []IntegrationPackage `json:"results"`
	} `json:"d"`
}

func (r *IPResponse) Print() {
	//	fmt.Println("ID Name")

	var text string
	for _, ip := range r.D.Results {
		text += fmt.Sprintf(
			"%s %s\n",
			ip.ID, ip.Name)
	}

	//	println(text)

	var responsePrinter IPResponsePrinter

	for _, ip := range r.D.Results {
		ipprinter := IPPrinter{
			ID:   ip.ID,
			Name: ip.Name,
			//			Description:     ip.Description,
			//			ShortText:       ip.ShortText,
			Version: ip.Version,
			Vendor:  ip.Vendor,
			//		PartnerContent:  ip.PartnerContent,
			UpdateAvailable: ip.UpdateAvailable,
			Mode:            ip.Mode,
			CreatedBy:       ip.CreatedBy,
		}
		responsePrinter.D.Results = append(responsePrinter.D.Results, ipprinter)

	}

	tableprinter.Print(os.Stdout, responsePrinter.D.Results)
}

type IPByIdResponse struct {
	D IntegrationPackage `json:"d"`
}

func (r *IPByIdResponse) Print() {
	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		panic("Could not Marshal IPByIdResponse")
	}
	fmt.Println(string(b))
}

type IntegrationFlow struct {
	Metadata    Metadata `json:"__metadata"`
	ID          string   `json:"Id"`
	Version     string   `json:"Version"`
	PackageID   string   `json:"PackageId"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Sender      string   `json:"Sender"`
	Receiver    string   `json:"Receiver"`
}

type Metadata struct {
	ID          string `json:"id"`
	URI         string `json:"uri"`
	Type        string `json:"type"`
	ContentType string `json:"content_type"`
	MediaSrc    string `json:"media_src"`
	EditMedia   string `json:"edit_media"`
}

type FlowsOfIPResponse struct {
	D struct {
		Results []IntegrationFlow `json:"results"`
	} `json:"d"`
}

type IPPrinter struct {
	ID   string `header:"Id"`
	Name string `header:"Name"`
	//	Description     string `header:"Description"`
	//	ShortText       string `header:"ShortText"`
	Version string `header:"Version"`
	Vendor  string `header:"Vendor"`
	//	PartnerContent  bool   `header:"PartnerContent"`
	UpdateAvailable bool   `header:"UpdateAvailable"`
	Mode            string `header:"Mode"`
	CreatedBy       string `header:"CreatedBy"`
}

type IPResponsePrinter struct {
	D struct {
		Results []IPPrinter
	}
}
