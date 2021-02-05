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
	"encoding/base64"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrmessages"
)

// GatewayClientEstablishment sends a GatewayClientEstablishmentRequest and processes a response.
// Return true if the establishment request was successful.
func (c *Comms) GatewayClientEstablishment(challenge [32]byte) (bool, error) {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(challenge)))
	base64.StdEncoding.Encode(b, challenge[:])

	request, err := fcrmessages.EncodeClientEstablishmentRequest(c.settings.ClientID(), string(b), c.settings.EstablishmentTTL())
	if err != nil {
		logging.Error("Error encoding Client Establishment Request: %+v", err)
		return false, err
	}

	if request.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(c.settings.RetrievalPrivateKey(), c.settings.RetrievalPrivateKeyVer(), msg)
	}) != nil {
		logging.Error("Error signing message for Client Establishment Request: %+v", err)
		return false, err
	}

	res := c.gatewayCall(request).Get("result").MustString()
	// TODO interpret the response.
	logging.Info("Response from server: %s", res)

	return true, nil
}

