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
	settings     settings.ClientSettings
	gateways     []ActiveGateway
	gatewaysLock sync.RWMutex
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info  register.GatewayRegister
	comms *gatewayapi.Comms
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
	gws, err := register.GetRegisteredGateways(g.settings.RegisterURL())
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
	}
	logging.Info("Register returned %d gateways", len(gws))

	// TODO this loop is where the managing of gateways that the client is using happens.
	logging.Info("Gateway Manager: GetGateways returned %d gateways", len(gws))
	for _, gw := range gws {
		gatewayID := gw.NodeID
		if g.gatewayInfoInvalid(&gw) {
			logging.Warn("Gateway registration information for gateway (%s) is invalid. Ignoring.", gatewayID)
			continue
		}

		logging.Info("Setting-up comms with: %+v", gw)
		comms, err := gatewayapi.NewGatewayAPIComms(&gw, &g.settings)
		if err != nil {
			panic(err)
		}

		// Try to do the establishment with the new gateway
		var challenge [32]byte
		fcrcrypto.GeneratePublicRandomBytes(challenge[:])
		comms.GatewayClientEstablishment(challenge)

		activeGateway := ActiveGateway{gw, comms}
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

	var aggregateOffers []cidoffer.CidGroupOffer
	for _, gw := range g.gateways {
		// TODO need to do nonce management
		// TODO need to do requests to all gateways in parallel, rather than serially
		offers, err := gw.comms.GatewayStdCIDDiscovery(contentID, 1)
		if err != nil {
			logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.info.NodeID, err)
		}
		// TODO: probably should remove duplicate offers at this point
		aggregateOffers = append(aggregateOffers, offers...)
	}
	return aggregateOffers, nil
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


// Validate the information coming from the Gateway.
func (g *GatewayManager) gatewayInfoInvalid(gateway *register.GatewayRegister) bool {
	// All of the fields must have a value in them.
	if gateway.NodeID == "" {
		logging.Warn("Gateway registration issue: NodeID not set")
		return true;
	}
	if gateway.Address == "" {
		logging.Warn("Gateway registration issue: Gateway IP address or domain name not set")
		return true;
	}
	if gateway.NetworkGatewayInfo == "" {
		logging.Warn("Gateway registration issue: Port for Gateway to Gateway communications not set")
		return true;
	}
	if gateway.NetworkProviderInfo == "" {
		logging.Warn("Gateway registration issue: Port for Provider to Gateway communications not set")
		return true;
	}
	if gateway.NetworkClientInfo == "" {
		logging.Warn("Gateway registration issue: Port for Client to Gateway communications not set")
		return true;
	}
	if gateway.NetworkAdminInfo == "" {
		logging.Warn("Gateway registration issue: Port for Admin to Gateway communications not set")
		return true;
	}
	if gateway.RegionCode == "" {
		logging.Warn("Gateway registration issue: Region Code not set")
		return true;
	}
	if gateway.RootSigningKey == "" {
		logging.Warn("Gateway registration issue: Root Signing Public Key not set")
		return true;
	}
	if gateway.SigingKey == ""	{
		logging.Warn("Gateway registration issue: Retrieval Signing Key not set")
		return true;
	}
	return false;
}