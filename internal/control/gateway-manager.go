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

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings settings.ClientSettings
	gateways []ActiveGateway
	gatewaysLock   sync.RWMutex

	// Registered Gateways
	RegisteredGateways []register.GatewayRegister
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info register.GatewayRegister
	comms 		*gatewayapi.Comms

}

// NewGatewayManager returns an initialised instance of the gateway manager.
func NewGatewayManager(settings settings.ClientSettings) *GatewayManager {
	g := GatewayManager{}
	g.settings = settings
	g.gatewayManagerRunner()
	return &g
}


// TODO this should be in a go routine and loop for ever.
func (g *GatewayManager) gatewayManagerRunner() {
	logging.Info("Gateway Manager: Management thread started")


	// Call this once each hour or maybe day.
	gateways, err := register.GetRegisteredGateways("http://register:8090")
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	g.RegisteredGateways = gateways

	// TODO this loop is where the managing of gateways that the client is using happens.
	logging.Info("Gateway Manager: GetGateways returned %d gateways", len(gateways))
	for _, gateway := range gateways {
		logging.Info("Setting-up comms with: %+v", gateway)
		comms, err := gatewayapi.NewGatewayAPIComms(&gateway, &g.settings)
		if err != nil {
			panic(err)
		} 

		// Try to do the establishment with the new gateway
		var challenge [32]byte
		fcrcrypto.GeneratePublicRandomBytes(challenge[:])
		comms.GatewayClientEstablishment(challenge)

		activeGateway := ActiveGateway{gateway, comms}
		g.gateways = append(g.gateways, activeGateway)
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


// FindOffersStandardDiscovery finds offers using the standard discovery mechanism.
func (g *GatewayManager) FindOffersStandardDiscovery(contentID *cid.ContentID) ([]cidoffer.CidGroupOffer, error) {
	if len(g.gateways) == 0 {
		return nil, fmt.Errorf("No gateways available")
	}
	return nil, nil

	// TODO need to do nonce management
//	for gateway := range g.gateways {
// TODO		gateway.comms.GatewayStdCIDDiscovery(contentID, 1) 

//	}

}

// GetConnectedGateways returns the list of domain names of gateways that the client 
// is currently connected to.
func (g *GatewayManager) GetConnectedGateways() []string {
	urls := make([]string, len(g.gateways))
	for i, gateway := range g.gateways {
		urls[i] = gateway.comms.ApiURL
	}
	return urls
}


// Shutdown stops go routines and closes sockets. This should be called as part 
// of the graceful library shutdown
func (g *GatewayManager) Shutdown() {
	// TODO
}
