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

import (
	"log"
	"net"
//	"os"
//	"time"
	"errors"
)

// GatewayPing sends "ping" message to gateway
func (c *FilecoinRetrievalClient) gatewayPing(server string) (bool, error) {
	if len(server) == 0 {
		errors.New("Error: Cannot ping empty servername")
	} else {
		log.Println("Attempting to ping \"" + server + "\"")
	}

	ra, err := net.ResolveIPAddr("ip4", server)
	if err != nil {
		errors.New("Error: Cannot ping empty servername")
	}
	log.Printf("Resolved %s as %s\n", server, ra.String())



	args := make(map[string]interface{})
	// TODO have a random challenge
	args["challenge"] = "123456789"
	args["ttl"] = "100"

	res := call("establishment", args).Get("result").MustString()
	log.Printf("Response from server: %s\n", res)

	return true, nil
}

