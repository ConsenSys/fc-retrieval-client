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
	GatewayID *big.Int
	Hostname string
	Location *LocationInfo
	// TODO public key?
}

// Information about the location of a gateway
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
		g = GatewayRegistrationContract
		createDummyData()
		singleInstance = &g
	})
	return singleInstance
}

func (g *GatewayRegistrationContract) createDummyData() {
	l = LocationInfo{RegionCode: "A", CountryCode: "AU", SubDivisionCode: "AU-QLD"}
	gi = GatewayInformation{GatewayID: big.NewInt(1), Hostname: "gateway", Location: &l)
	append(g.gateways, gi)
}

// GetGateways returns gateways anywhere in the world
func (g *GatewayRegistrationContract) GetGateways(maxToReturn int32) *[]GatewayInformation {
	// TODO handle num == 0

	return g.gateways
}

// GetGateways returns gateways that match the region code
func (g *GatewayRegistrationContract) GetGateways(maxToReturn int32, regionCode string) *[]GatewayInformation {
	// TODO handle num == 0
	if regionCode != "A" {
		return noGateways
	}

	return g.gateways
}

// GetGateways returns gateways that match the region code and the country code
func (g *GatewayRegistrationContract) GetGateways(maxToReturn int32, regionCode string, countryCode string) *[]GatewayInformation {
	// TODO handle num == 0
	if regionCode != "A" || countryCode != "AU" {
		return noGateways
	}

	return g.gateways
}

// GetGateways returns gateways that match the precise location
func (g *GatewayRegistrationContract) GetGateways(maxToReturn int32, regionCode string, countryCode string, subdivisionCode string) *[]GatewayInformation {
	// TODO handle num == 0
	if regionCode != "A" || countryCode != "AU" || subdivisionCode != "AU-QLD" {
		return noGateways
	}

	return g.gateways
}
