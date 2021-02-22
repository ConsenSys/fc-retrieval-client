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
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

	"github.com/ConsenSys/fc-retrieval-client/internal/gatewayapi"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// GatewayManager managers the pool of gateways and the connections to them.
type GatewayManager struct {
	settings     settings.ClientSettings
	gateways     []ActiveGateway
	gatewaysLock sync.RWMutex

	// List of gateway to use. A client may request a node be added to this list e
	gatewaysToUse []*nodeid.NodeID

	done 				 chan bool
	ticker       *time.Ticker
}

// ActiveGateway contains information for a single gateway
type ActiveGateway struct {
	info  register.GatewayRegister
	comms *gatewayapi.Comms
	nodeID *nodeid.NodeID
}

// NewGatewayManager returns an initialised instance of the gateway manager.
func NewGatewayManager(settings settings.ClientSettings) *GatewayManager {
	g := GatewayManager{}
	g.settings = settings
	g.gatewaysToUse = make([]*nodeid.NodeID, 0)
	g.gatewayManagerRunner()
	return &g
}


// Get the latest gateway information now.
func (g *GatewayManager) requestUpdate() {

}



// gatewayManagerRunner gets the latest gateway information. 
func (g *GatewayManager) gatewayManagerRunner() {
	logging.Info("Gateway Manager: Management thread started")

	g.ticker = time.NewTicker(30 * time.Second)
	g.done = make(chan bool)

	go func() {
		for {
				select {
				case <-g.done:
						return
				case <-g.ticker.C:
						g.getLatestGatewayInfo()
				}
		}
	}()
}

// get the latest gateway information from the registry.
// Note that this is run inside a go routine. No e
func (g *GatewayManager) getLatestGatewayInfo() {
	// Remove any gateways that are no longer in the list of gateways to use.
	gatewaysToRemove := make([]*nodeid.NodeID, 0)
	for _, gwInfo := range g.gateways {
		notFound := true
		for _, gwNodeID := range g.gatewaysToUse {
			if gwNodeID.ToString() == gwInfo.nodeID.ToString() {
				notFound = false
				break;
			}
		}
		if notFound {
			gatewaysToRemove = append(gatewaysToRemove, gwInfo.nodeID)
		}
	}
	for _, gwNodeID := range gatewaysToRemove {
		for i, gwInfo := range g.gateways {
			if gwNodeID.ToString() == gwInfo.nodeID.ToString() {
				g.gateways[i] = g.gateways[len(g.gateways)-1]
				g.gateways = g.gateways[:len(g.gateways)-1]
			}
		}
	}


	// Get the latest information from the register for exixting gateways
	// for _, gwInfo := range g.gateways {
	// 	// TODO
	// }


	// Add information for new gateways.
	for _, gwNodeID := range g.gatewaysToUse {
		found := false
		for _, gwInfo := range g.gateways {
			if gwNodeID.ToString() == gwInfo.nodeID.ToString() {
				found = true
				break;
			}
		}
		if !found {
			g.addGateway(gwNodeID)
		}
	}
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
	g.ticker.Stop()
	g.done <- true
}

func (g *GatewayManager) addGateway(nodeID *nodeid.NodeID) {
	// TODO add gateway by ID
	gws, err := register.GetRegisteredGateways(g.settings.RegisterURL())
	if err != nil {
		logging.Error("Unable to get registered gateways: %v", err)
		return
	}
	logging.Info("Register returned %d gateways", len(gws))
	if len(gws) == 0 {
		logging.Warn("Unable to get registered gateways: %v", err)
		return
	}
	if len(gws) != 1 {
		logging.Warn("Unexpectedly, multiple gateways returned: %d", len(gws))
		return
	}
	gw := gws[0]
	gatewayID, err := nodeid.NewNodeIDFromString(gw.NodeID)

	if !g.validateGatewayInfo(&gw) {
		logging.Warn("Gateway registration information for gateway (%s) is invalid. Ignoring.", gatewayID)
		return
	}

	logging.Info("Setting-up comms with: %+v", gw)
	comms, err := gatewayapi.NewGatewayAPIComms(&gw, &g.settings)
	if err != nil {
		logging.Error("Error encountered which contacting gateway (%s): %+v", gatewayID, err)
		return
	}

	// TODO this should only be done for new gateways.
	// Try to do the establishment with the new gateway
	var challenge [32]byte
	fcrcrypto.GeneratePublicRandomBytes(challenge[:])
	comms.GatewayClientEstablishment(challenge)
	if err != nil {
		logging.Warn("Error processing node id: %+v", err)
		return
	}

	activeGateway := ActiveGateway{gw, comms, gatewayID}
	g.gateways = append(g.gateways, activeGateway)

	if len(g.gateways) > 0 {
		logging.Info("Gateway Manager using %d gateways", len(g.gateways))
	} else {
		logging.Warn("No gateways available")
	}
}

