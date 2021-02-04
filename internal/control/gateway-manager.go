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
	"fmt"
	"sync"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/internal/contracts"
	"github.com/ConsenSys/fc-retrieval-client/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings settings.ClientSettings
	// gatewayRegistrationContract *contracts.GatewayRegistrationContract
	gateways     []ActiveGateway
	gatewaysLock sync.RWMutex
	// Registered Gateways
	RegisteredGateways []register.GatewayRegister
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info  contracts.GatewayInformation
	comms *gatewayapi.Comms
}

var doOnce sync.Once
var singleInstance *GatewayManager

// GetGatewayManager returns the single instance of the gateway manager.
// The settings parameter must be used with the first call to this function.
// After that, the settings parameter is ignored.
func GetGatewayManager(settings ...settings.ClientSettings) *GatewayManager {
	doOnce.Do(func() {
		if len(settings) != 1 {
			// TODO replace with ErrorAndPanic once available
			logging.ErrorAndPanic("Unexpected number of parameter passed to first call of GetGatewayManager")
		}
		startGatewayManager(settings[0])
	})
	return singleInstance
}

func startGatewayManager(settings settings.ClientSettings) {
	g := GatewayManager{}
	g.settings = settings

	singleInstance = &g
	g.gatewayManagerRunner()
}

func (g *GatewayManager) gatewayManagerRunner() {
	logging.Info("Gateway Manager: Management thread started")

	// Call this once each hour or maybe day.
	gateways, err := register.GetRegisteredGateways("http://fc-retrieval-register:8090")
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	g.RegisteredGateways = gateways

	// gatewayInfo := g.gatewayRegistrationContract.GetGateways(10)
	logging.Info("Gateway Manager: GetGateways returned %d gateways", len(gateways))
	for _, gateway := range gateways {

		fmt.Printf("Gateway ========> %+v\n", gateway)
		// comms, err := gatewayapi.NewGatewayAPIComms(info.Hostname, g.settings.ClientID())
		// if err != nil {
		// 	panic(err)
		// }

		// Try to do the establishment with the new gateway
		// var challenge [32]byte
		// fcrcrypto.GenerateRandomBytes(challenge[:])
		// comms.GatewayClientEstablishment(g.settings.EstablishmentTTL(), challenge)

		// activeGateway := ActiveGateway{info, comms}
		// g.gateways = append(g.gateways, activeGateway)
	}

	logging.Info("Gateway Manager using %d gateways", len(g.gateways))
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
