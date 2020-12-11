package contracts
/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

 import (
	 "sync"
 )

 
/*
* This file is a wrapper for the Gateway Registration Contract. At present it returns
* hard coded values.
*/


// GatewayRegistrationContract provides a wrapper for the real Gateway Registration Contract
type GatewayRegistrationContract struct {
	gateways []GatewayInformation
}

// GatewayInformation holds information about a single gateway.
type GatewayInformation struct {
	GatewayID [32]byte
	Hostname string
	Location *LocationInfo
	// TODO public key?
}

// LocationInfo contains information about the location of a gateway
type LocationInfo struct {
	RegionCode string
	CountryCode string
	SubDivisionCode string
}

var doOnce sync.Once
var singleInstance *GatewayRegistrationContract
var noGateways []GatewayInformation

// GetGatewayRegistrationContract gets the single instance of the wrapper for the 
// gateway registration contract.
func GetGatewayRegistrationContract() *GatewayRegistrationContract {
    doOnce.Do(func() {
		g := GatewayRegistrationContract{}
		g.createDummyData()
		singleInstance = &g
	})
	return singleInstance
}

func (g *GatewayRegistrationContract) createDummyData() {
	l := LocationInfo{RegionCode: "A", CountryCode: "AU", SubDivisionCode: "AU-QLD"}
	var dummyGatewayID [32]byte
	dummyGatewayID[0] = 0x12
	gi := GatewayInformation{GatewayID: dummyGatewayID, Hostname: "gateway", Location: &l}

	g.gateways = append(g.gateways, gi)
}

// GetGateways returns gateway information based on the locaiton parameters
func (g *GatewayRegistrationContract) GetGateways(maxToReturn int32, loc ...string) []GatewayInformation {
	len := len(loc)
	switch len {
	case 0:
		// TODO only return up to maxToReturn
		return g.gateways
	case 1:
		regionCode := loc[0]
		// TODO only return up to maxToReturn
		// TODO this assumes the hard coded data
		if regionCode != "A" {
			return noGateways
		}
		return g.gateways
	case 2:
		regionCode := loc[0]
		countryCode := loc[1]
		// TODO only return up to maxToReturn
		// TODO this assumes the hard coded data
		if regionCode != "A" || countryCode != "AU" {
			return noGateways
		}
		return g.gateways
	case 3:
		regionCode := loc[0]
		countryCode := loc[1]
		subdivisionCode := loc[2]
		// TODO only return up to maxToReturn
		// TODO this assumes the hard coded data
		if regionCode != "A" || countryCode != "AU" || subdivisionCode != "AU-QLD" {
			return noGateways
		}
		return g.gateways
	default:
		panic("Invalid number of parameters to contract.GetGateways")
	}
}
