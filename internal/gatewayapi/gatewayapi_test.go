package gatewayapi

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-client/internal/contracts"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)

func TestSigning(t *testing.T) {
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	retirevalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		panic(err)
	}
	retrievalPrivateKeyVer := fcrcrypto.DecodeKeyVersion(1)


	s := settings.CreateSettings()
	s.SetBlockchainPrivateKey(blockchainPrivateKey)
	s.SetRetrievalPrivateKey(retirevalPrivateKey, retrievalPrivateKeyVer);
	settings := s.Build()

	gAPI, err := NewGatewayAPIComms(createDummyGatewayInformation(), settings)
	if err != nil {
		panic(err)
	}
	// Switch out the gateway key for the client key for this test so that the 
	// client message can be verified using the client public key.
	gAPI.gatewayPubKey = retirevalPrivateKey

	msg := messages.ClientEstablishmentRequest{}
	msg.Challenge = "1234567890abcdef1234567890"
	method := int32(messages.ClientEstablishmentRequestType)
	gAPI.addCommonFieldsAndSign(method, &msg.ClientCommonRequestFields, msg);
	logging.Test("message: %+v", msg)

	// TODO verify
	signature := msg.ClientCommonRequestFields.Signature
	msg.ClientCommonRequestFields.Signature = ""
	verified := gAPI.verifyMessage(signature, msg)
	assert.Equal(t, true, verified)
}



func createDummyGatewayInformation() *contracts.GatewayInformation{
	l := contracts.LocationInfo{RegionCode: "A", CountryCode: "AU", SubDivisionCode: "AU-QLD"}
	var dummyGatewayID [32]byte
	dummyGatewayID[0] = 0x12
	gatewayKeyPair, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		logging.ErrorAndPanic("Error: %s", err)
	}
	encodedPubKey, err := gatewayKeyPair.EncodePublicKey()
	if err != nil {
		logging.ErrorAndPanic("Error: %s", err)
	}
	gatewayPublicKey, err := fcrcrypto.DecodePublicKey(encodedPubKey)

	gatewayPublicKeyVer := fcrcrypto.InitialKeyVersion()
	
	gi := contracts.GatewayInformation{
		GatewayID: dummyGatewayID, 
		Hostname: "localhost",   // use a host name that will definitely resolve
		Location: &l, 
		GatewayRetrievalPublicKey: gatewayPublicKey,
		GatewayRetrievalPublicKeyVersion: gatewayPublicKeyVer,
	}
	return &gi
}