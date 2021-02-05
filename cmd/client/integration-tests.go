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

	_ "github.com/joho/godotenv/autoload"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-client/config"
	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
)

func main() {
	conf := config.NewConfig()
	logging.Init(conf)
	logging.Debug("Using arg client-id=%v", conf.GetString("CLIENT_ID"))
	logging.Debug("Using arg ttl=%v", conf.GetInt("ESTABLISHMENT_TTL"))
	logging.Debug("Using arg log-level=%v", conf.GetString("LOG_LEVEL"))
	logging.Debug("Using arg log-target=%v", conf.GetString("LOG_TARGET"))

	// TODO switch this to logging.Test when available
	logging.Info("Integration Test: Start")
	integrationTests()
	logging.Info("Integration Test: End")
}

func integrationTests() {
	// TODO switch this to logging.Test when available
	logging.Info(" Wait two seconds for the gateway to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	var pieceCIDToFind [32]byte


	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	client := fcrclient.InitFilecoinRetrievalClient(*conf)
	offers := client.FindBestOffers(pieceCIDToFind, 1000, 1000)
	// TODO switch this to logging.Test when available
	logging.Info("Offers: %+v\n", offers)
	client.Shutdown()
}
