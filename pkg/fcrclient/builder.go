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
	"math/big"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder struct {
	logLevel         string
	logTarget        string
	logServiceName   string
	establishmentTTL int64
	clientID         *nodeid.NodeID
	registerURL      string

	blockchainPrivateKey *fcrcrypto.KeyPair

	retrievalPrivateKey    *fcrcrypto.KeyPair
	retrievalPrivateKeyVer *fcrcrypto.KeyVersion

	walletPrivateKey string
	lotusAP          string
	lotusAuthToken   string
	searchPrice      *big.Int
	offerPrice       *big.Int
	topUpAmount      *big.Int
}

// CreateSettings creates an object with the default settings.
func CreateSettings() *SettingsBuilder {
	f := SettingsBuilder{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.logServiceName = defaultLogServiceName
	f.establishmentTTL = defaultEstablishmentTTL
	f.registerURL = defaultRegisterURL
	f.searchPrice = big.NewInt(defaultSearchPrice)
	f.offerPrice = big.NewInt(defaultOfferPrice)
	f.topUpAmount = big.NewInt(defaultTopUpAmount)
	return &f
}

// SetLogging sets the log level and target.
func (f *SettingsBuilder) SetLogging(logLevel string, logTarget string, logServiceName string) {
	f.logLevel = logLevel
	f.logTarget = logTarget
	f.logServiceName = logServiceName
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
func (f *SettingsBuilder) SetEstablishmentTTL(ttl int64) {
	f.establishmentTTL = ttl
}

// SetRegisterURL sets the register URL.
func (f *SettingsBuilder) SetRegisterURL(url string) {
	f.registerURL = url
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *SettingsBuilder) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f *SettingsBuilder) SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.retrievalPrivateKey = rPkey
	f.retrievalPrivateKeyVer = ver
}

func (f *SettingsBuilder) SetWalletPrivateKey(walletPrivateKey string) {
	f.walletPrivateKey = walletPrivateKey
}

func (f *SettingsBuilder) SetLotusAP(lotusAP string) {
	f.lotusAP = lotusAP
}

func (f *SettingsBuilder) SetLotusAuthToken(lotusAuthToken string) {
	f.lotusAuthToken = lotusAuthToken
}

func (f *SettingsBuilder) SetSearchPrice(searchPrice *big.Int) {
	f.searchPrice = searchPrice
}

func (f *SettingsBuilder) SetOfferPrice(offerPrice *big.Int) {
	f.offerPrice = offerPrice
}

func (f *SettingsBuilder) SetTopUpAmount(topUpAmount *big.Int) {
	f.topUpAmount = topUpAmount
}

// Build creates a settings object and initialises the logging system.
func (f *SettingsBuilder) Build() *ClientSettings {

	logging.Init1(f.logLevel, f.logTarget, f.logServiceName)

	g := ClientSettings{}
	g.establishmentTTL = f.establishmentTTL
	g.registerURL = f.registerURL

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.clientID == nil {
		logging.Info("Settings: No Client ID set. Generating random client ID")
		// TODO replace once NewRandomNodeID becomes available.
		g.clientID = nodeid.NewRandomNodeID()
	} else {
		g.clientID = f.clientID
	}

	if f.retrievalPrivateKey == nil {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s", err)
		}
		g.retrievalPrivateKey = pKey
		g.retrievalPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.retrievalPrivateKey = f.retrievalPrivateKey
		g.retrievalPrivateKeyVer = f.retrievalPrivateKeyVer
	}

	// TODO add validations
	g.walletPrivateKey = f.walletPrivateKey
	g.lotusAP = f.lotusAP
	g.lotusAuthToken = f.lotusAuthToken
	g.searchPrice = f.searchPrice
	g.offerPrice = f.offerPrice
	g.topUpAmount = f.topUpAmount

	return &g
}
