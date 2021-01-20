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

	//	"github.com/stretchr/testify/assert"

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

	gAPI, err := NewGatewayAPIComms("1.2.3.4", settings)
	if err != nil {
		panic(err)
	}

	msg := messages.ClientEstablishmentRequest{}
	msg.Challenge = "1234567890abcdef1234567890"
	method := int32(messages.ClientEstablishmentRequestType)
	gAPI.addCommonFieldsAndSign(method, msg.ClientCommonRequestFields);
	logging.Error("message type: %+v", msg.MessageType)

	// TODO verify
}