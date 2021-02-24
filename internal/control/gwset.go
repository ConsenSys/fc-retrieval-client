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
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"

//	"github.com/ConsenSys/fc-retrieval-client/internal/gatewayapi"
)


// FindGateways find gateways located near too the specified location. Use AddGateways
// to use these gateways.
func (g *GatewayManager) FindGateways(location []string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	// Determine gateways to use. For the moment, this is just "use all of them"
	// TODO: This will have to become, use gateways that this client has FIL registered with.
	gws, err := register.GetRegisteredGateways(g.settings.RegisterURL())
	if err != nil {
		return nil, err
	}

	var gatewayNodes []*nodeid.NodeID
	for _, registrationInfo := range gws {
		nodeid, err := nodeid.NewNodeIDFromString(registrationInfo.NodeID)
		if err != nil {
			logging.Warn("FindGateways: Ignoring invalid node id: %s", registrationInfo.NodeID)
			continue
		}
		if !g.validateGatewayInfo(&registrationInfo) {
			continue
		}

		gatewayNodes = append(gatewayNodes, nodeid)
	}

	return gatewayNodes, nil
}

// AddGateways adds one or more gateways to use.
// Returns the number of gateways added.
func (g *GatewayManager) AddGateways(gwNodeIDs []*nodeid.NodeID) int {
	g.gatewaysToUseLock.RLock()
	defer g.gatewaysToUseLock.RUnlock()

	numAdded := 0
	for _, gwToUseID := range gwNodeIDs {
		add := true
		for _, gwID := range g.gatewaysToUse {
			if gwID.ToString() == gwToUseID.ToString() {
				add = false
			}
		}
		if add {
			g.gatewaysToUse = append(g.gatewaysToUse, gwToUseID)
			numAdded++
		}
	}
	g.requestUpdate()

	return numAdded
}

// RemoveGateways removes one or more gateways from the list of Gateways to use.
func (g *GatewayManager) RemoveGateways(gwNodeIDs []*nodeid.NodeID) int {
	g.gatewaysToUseLock.RLock()
	defer g.gatewaysToUseLock.RUnlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		for i, gwID := range g.gatewaysToUse {
			if gwID.ToString() == gwToRemoveID.ToString() {
				// Remove
				g.gatewaysToUse[i] = g.gatewaysToUse[len(g.gatewaysToUse)-1]
				g.gatewaysToUse = g.gatewaysToUse[:len(g.gatewaysToUse)-1]
				numRemoved++
				break
			}
		}
	}
	g.requestUpdate()

	return numRemoved
}

// RemoveAllGateways removes all gateways from the list of Gateways to use.
func (g *GatewayManager) RemoveAllGateways() int {
	g.gatewaysToUseLock.RLock()
	defer g.gatewaysToUseLock.RUnlock()

	numRemoved := len(g.gatewaysToUse)
	g.gatewaysToUse = g.gatewaysToUse[:0]

	g.requestUpdate()
	return numRemoved
}

// GetGateways returns the list of gateways that are being used.
func (g *GatewayManager) GetGateways() []*nodeid.NodeID {
	return g.gatewaysToUse
}
