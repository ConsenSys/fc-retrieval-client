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
	"time"

	// "github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
)

func main() {
	logging.Error("Integration Test: Start")
	integrationTests()
	logging.Error("Integration Test: End")
}

func integrationTests() {
	logging.Error("Wait two seconds for the gateway to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	var pieceCIDToFind [32]byte

	// TEMP REMOVE
	// blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	// if err != nil {
	// 	panic(err)
	// }
	// blockchainPrivateKeyAlg := fcrcrypto.DecodeSigAlg(fcrcrypto.SigAlgEcdsaP256Sha512_256)

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	// confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey, blockchainPrivateKeyAlg)
	conf := confBuilder.Build()

	client := fcrclient.InitFilecoinRetrievalClient(*conf)
	offers := client.FindBestOffers(pieceCIDToFind, 1000, 1000)
	logging.Error("Offers: %+v\n", offers)
	client.Shutdown()
}
