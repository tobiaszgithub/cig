package model

import (
	"fmt"
	"os"

	"github.com/lensesio/tableprinter"
)

type IntegrationPackage struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type IPResponse2 struct {
	D struct {
		Results []struct {
			ID   string `json:"Id"`
			Name string `json:"Name"`
		} `json:"results"`
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
		}
		responsePrinter.D.Results = append(responsePrinter.D.Results, ipprinter)

	}

	tableprinter.Print(os.Stdout, responsePrinter.D.Results)
}

type IPResponse struct {
	D struct {
		Results []struct {
			ID                string      `json:"Id"`
			Name              string      `json:"Name"`
			Description       string      `json:"Description"`
			ShortText         string      `json:"ShortText"`
			Version           string      `json:"Version"`
			Vendor            string      `json:"Vendor"`
			PartnerContent    bool        `json:"PartnerContent"`
			UpdateAvailable   bool        `json:"UpdateAvailable"`
			Mode              string      `json:"Mode"`
			SupportedPlatform string      `json:"SupportedPlatform"`
			ModifiedBy        string      `json:"ModifiedBy"`
			CreationDate      string      `json:"CreationDate"`
			ModifiedDate      string      `json:"ModifiedDate"`
			CreatedBy         string      `json:"CreatedBy"`
			Products          string      `json:"Products"`
			Keywords          string      `json:"Keywords"`
			Countries         string      `json:"Countries"`
			Industries        string      `json:"Industries"`
			LineOfBusiness    string      `json:"LineOfBusiness"`
			PackageContent    interface{} `json:"PackageContent"`
		} `json:"results"`
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
}

type IPResponsePrinter struct {
	D struct {
		Results []IPPrinter
	}
}
