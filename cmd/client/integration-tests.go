package main

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
	"log"
	"time"
	"math/big"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
)

func main() {
	log.Println("Integration Test: Start")
	integrationTests()
	log.Println("Integration Test: End")
}

func integrationTests() {
	log.Println(" Wait two seconds for the gateway to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	var pieceCIDToFind [32]byte

	// TODO remove this and handle node ID
	nodeID, err := nodeid.NewNodeID(big.NewInt(7))
	if err != nil {
		panic(err)
	}
	settings := fcrclient.FilecoinRetrievalClientSettings{MaxEstablishmentTTL: 100, Verbose: true, NodeID: nodeID}
	client := fcrclient.InitFilecoinRetrievalClient(&settings)
	offers := client.FindBestOffers(pieceCIDToFind, 1000, 1000)
	log.Printf("Offers: %+v\n", offers)
	client.Shutdown()
}
