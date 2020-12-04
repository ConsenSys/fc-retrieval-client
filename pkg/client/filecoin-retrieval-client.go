package client

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

// FilecoinRetrievalClient holds information about the interaction of 
// the Filecoin Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	// TODO have a list of gateway objects of all the current gateways being interacted with
}

// NewFilecoinRetrievalClient creates a Filecoin Retrieval Client
// TODO: should this return a reference to a singleton?
func NewFilecoinRetrievalClient() *FilecoinRetrievalClient {
	var c = FilecoinRetrievalClient{}
	return &c
}

// AddGateway should probably not be exposed as a public function.
// This should happen automatically.
func (c *FilecoinRetrievalClient) AddGateway(hostname string) {
	// For the moment, just ping the gateway
	c.gatewayPing(hostname)
}