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
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// ClientSettings holds the library configuration
type ClientSettings struct {
	establishmentTTL int64
	registerURL      string
	clientID         *nodeid.NodeID

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

// WalletPrivateKey returns the wallet private key
func (c ClientSettings) WalletPrivateKey() string {
	return c.walletPrivateKey
}

// LotusAP returns the lotusAP
func (c ClientSettings) LotusAP() string {
	return c.lotusAP
}

// LotusAuthToken returns the lotus authorization token
func (c ClientSettings) LotusAuthToken() string {
	return c.lotusAuthToken
}

// SearchPrice returns the search price
func (c ClientSettings) SearchPrice() *big.Int {
	return c.searchPrice
}

// OfferPrice returns offer price
func (c ClientSettings) OfferPrice() *big.Int {
	return c.offerPrice
}

// TopUpAmount returns the  top up amount
func (c ClientSettings) TopUpAmount() *big.Int {
	return c.topUpAmount
}

// EstablishmentTTL returns the establishmentTTL
func (c ClientSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// RegisterURL returns the register URL
func (c ClientSettings) RegisterURL() string {
	return c.registerURL
}

// ClientID returns the ClientID
func (c ClientSettings) ClientID() *nodeid.NodeID {
	return c.clientID
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c ClientSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// RetrievalPrivateKey returns the RetrievalPrivateKey
func (c ClientSettings) RetrievalPrivateKey() *fcrcrypto.KeyPair {
	return c.retrievalPrivateKey
}

// RetrievalPrivateKeyVer returns the RetrievalPrivateKeyVer
func (c ClientSettings) RetrievalPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.retrievalPrivateKeyVer
}
