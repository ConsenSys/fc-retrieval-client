package fcrclient

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
	"encoding/hex"

	"github.com/ConsenSys/fc-retrieval-client/internal/control"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
)

// FilecoinRetrievalClient holds information about the interaction of
// the Filecoin Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	gatewayManager *control.GatewayManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

// NewFilecoinRetrievalClient initialise the Filecoin Retreival Client library
func NewFilecoinRetrievalClient(conf Settings) (*FilecoinRetrievalClient, error) {
	logging.Info("Filecoin Retrieval Client started")
	var c = FilecoinRetrievalClient{}
	clientSettings := conf.(*settings.ClientSettings)
	c.gatewayManager = control.GetGatewayManager(*clientSettings)
	return &c, nil

}


// FindBestOffers locates offsers for supplying the content associated with the pieceCID
func (c *FilecoinRetrievalClient) FindBestOffers(pieceCID [32]byte, maxPrice int64, maxExpectedLatency int64) ([]cidoffer.CidGroupOffer){
	var hexDumpPieceCID string
	if logging.InfoEnabled() {
		hexDumpPieceCID = hex.Dump(pieceCID[:])
		logging.Info("Filecoin Retrieval Client: FindBestOffers(pieceCID: %s, maxPrice: %d, maxExpectedLatency: %d", 
			hexDumpPieceCID, maxPrice, maxExpectedLatency)
	}

	// TODO
	logging.Info("Filecoin Retrieval Client: FindBestOffers(pieceCID: %s) returning no offers", hexDumpPieceCID)
	return nil
}

// Shutdown releases all resources used by the library
func (c *FilecoinRetrievalClient) Shutdown() {
	logging.Info("Filecoin Retrieval Client shutting down")
	c.gatewayManager.Shutdown()
}