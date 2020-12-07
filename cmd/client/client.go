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
	"os"
	"time"

	"github.com/ConsenSys/fc-retrieval-client/pkg/client"
)

var (
	servername = "localhost"
	pingserver = servername
)

func main() {
	log.Println("Integration Test: Start")
	integrationTests()
	log.Println("Integration Test: End")
}

func integrationTests() {
	log.Println(" Wait one second for the gateway to deploy and be ready for requests")
	time.Sleep(time.Second)

	//ping()
	addGateway()
}

func addGateway() {
	log.Println("Test: addGateway")
	fc := client.NewFilecoinRetrievalClient()
	fc.AddGateway("gateway")
}

func ping() {
	if len(os.Args[1:]) > 0 {
		pingserver = os.Args[1]
	}
	client.Ping(pingserver)
}
