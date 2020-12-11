package control

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

	"github.com/ConsenSys/fc-retrieval-client/internal/contracts"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	gatewayRegistrationContract *contracts.GatewayRegistrationContract 
	gateways []ActiveGateway
	gatewaysLock   sync.RWMutex
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info contracts.GatewayInformation

}

// GatewayManagerSettings is used to communicate the settings to be used by the 
// Gateway Manager.
type GatewayManagerSettings struct {
	MaxEstablishmentTTL int64
}

var doOnce sync.Once
var singleInstance *GatewayManager

// GetGatewayManager returns the single instance of the gateway manager.
// The settings parameter must be used with the first call to this function.
// After that, the settings parameter is ignored.
func GetGatewayManager(settings ...*GatewayManagerSettings) *GatewayManager {
    doOnce.Do(func() {
		if len(settings) != 1 {
			panic("Unexpected number of parameter passed to first call of GetGatewayManager")
		}
		startGatewayManager(settings[0])
	})
	return singleInstance
}

func startGatewayManager(settings *GatewayManagerSettings) {
	g := GatewayManager{}
	singleInstance = &g

//	errChan := make(chan error, 1)
	go g.gatewayManagerRunner()

	// TODO what should be done with error that is returned possibly in the future?
	// TODO would it be better just to have gatewayManagerRunner panic after emitting a log?
}

func (g *GatewayManager) gatewayManagerRunner() {
	contracts.GetGatewayRegistrationContract() 

	

	// TODO read gateway information from the smart contracts
}

// BlockGateway adds a host to disallowed list of gateways
func (g *GatewayManager) BlockGateway(hostName string) {
	// TODO
}


// UnblockGateway add a host to allowed list of gateways
func (g *GatewayManager) UnblockGateway(hostName string) {
	// TODO

}

// Shutdown stops go routines and closes sockets. This should be called as part 
// of the graceful library shutdown
func (g *GatewayManager) Shutdown() {
	// TODO
}
