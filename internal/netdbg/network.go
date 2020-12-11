package netdbg

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
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

// Ping pings a remote server via ICMP
func Ping(pingserver string) bool {

	if len(pingserver) == 0 {
		log.Println("Error: Cannot ping empty servername")
		os.Exit(1)
	} else {
		log.Println("Attempting to ping " + pingserver)
	}

	p := fastping.NewPinger()
	p.Network("udp")
	ra, err := net.ResolveIPAddr("ip4:icmp", pingserver)
	if err != nil {
		log.Println(err)
		return false
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		log.Println("finish")
	}
	err = p.Run()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
