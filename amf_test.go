// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

/*
 * AMF Unit Testcases
 *
 */
package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/omec-project/amf/factory"
	protos "github.com/omec-project/config5g/proto/sdcoreConfig"
)

//var AMF = &service.AMF{}

func init() {
	factory.InitConfigFactory("amfTest/amfcfg.yaml")
}

func GetNetworkSliceConfig() *protos.NetworkSliceResponse {
	var rsp protos.NetworkSliceResponse

	rsp.NetworkSlice = make([]*protos.NetworkSlice, 0)

	ns := protos.NetworkSlice{}
	slice := protos.NSSAI{Sst: "1", Sd: "010203"}
	ns.Nssai = &slice

	site := protos.SiteInfo{SiteName: "siteOne", Gnb: make([]*protos.GNodeB, 0), Plmn: new(protos.PlmnId)}
	gNb := protos.GNodeB{Name: "gnb", Tac: 1}
	site.Gnb = append(site.Gnb, &gNb)
	site.Plmn.Mcc = "208"
	site.Plmn.Mnc = "93"
	ns.Site = &site

	rsp.NetworkSlice = append(rsp.NetworkSlice, &ns)
	return &rsp
}

func TestInitialConfig(t *testing.T) {
	factory.AmfConfig.Configuration.PlmnSupportList = nil
	factory.AmfConfig.Configuration.ServedGumaiList = nil
	factory.AmfConfig.Configuration.SupportTAIList = nil
	var Rsp chan *protos.NetworkSliceResponse
	Rsp = make(chan *protos.NetworkSliceResponse)
	go func() {
		Rsp <- GetNetworkSliceConfig()
	}()
	go func() {
		AMF.UpdateConfig(Rsp)
	}()

	time.Sleep(2 * time.Second)
	if factory.AmfConfig.Configuration.PlmnSupportList != nil &&
		factory.AmfConfig.Configuration.ServedGumaiList != nil &&
		factory.AmfConfig.Configuration.SupportTAIList != nil {
		fmt.Printf("test passed")
	} else {
		t.Errorf("test failed")
	}
}

// data in JSON format which
// is to be decoded
var Data = []byte(`{
	"NetworkSlice": [
		{
		 "Name": "siteOne",
		 "Nssai": {"Sst": "1", "Sd": "010203"},
		 "Site": {
			"SiteName": "siteOne",
			"Gnb": [
				{"Name": "gnb1", "Tac": 1}, 
				{"Name": "gnb2", "Tac": 2}
			],
			"Plmn": {"mcc": "208", "mnc": "93"}
		  }
		}
		]}`)

func TestUpdateConfig(t *testing.T) {
	var nrp protos.NetworkSliceResponse
	err := json.Unmarshal(Data, &nrp)
	if err != nil {
		panic(err)
	}
	var Rsp chan *protos.NetworkSliceResponse
	Rsp = make(chan *protos.NetworkSliceResponse)
	go func() {
		Rsp <- &nrp
	}()
	go func() {
		AMF.UpdateConfig(Rsp)
	}()

	time.Sleep(2 * time.Second)
	if factory.AmfConfig.Configuration.SupportTAIList != nil &&
		len(factory.AmfConfig.Configuration.SupportTAIList) == 2 {
		fmt.Printf("test passed")
	} else {
		t.Errorf("test failed")
	}
}
