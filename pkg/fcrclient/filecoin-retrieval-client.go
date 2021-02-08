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
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"

)


// FilecoinRetrievalClient holds information about the interaction of 
// the Filecoin Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	gatewayManager *control.GatewayManager
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

var singleInstance *FilecoinRetrievalClient
var initialised = false

// InitFilecoinRetrievalClient initialise the Filecoin Retreival Client library
func InitFilecoinRetrievalClient(settings Settings) (*FilecoinRetrievalClient, error) {
	if initialised {
		panic("Attempt to init Filecoin Retrieval Client a second time")
	}
	var c = FilecoinRetrievalClient{}
	c.startUp(settings)
	singleInstance = &c
	initialised = true
	return singleInstance, nil

}




// GetFilecoinRetrievalClient creates a Filecoin Retrieval Client
func GetFilecoinRetrievalClient() *FilecoinRetrievalClient {
	if !initialised {
		panic("Filecoin Retrieval Client not initialised")
	}

	return singleInstance
}

func (c *FilecoinRetrievalClient) startUp(conf Settings) {
	logging.Info("Filecoin Retrieval Client started")
	clientSettings := conf.(*settings.ClientSettings)
	c.gatewayManager = control.GetGatewayManager(*clientSettings)
}




// FindBestOffers locates offsers for supplying the content associated with the pieceCID
func (c *FilecoinRetrievalClient) FindBestOffers(pieceCID [32]byte, maxPrice int64, maxExpectedLatency int64) ([]PieceCIDOffer){
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