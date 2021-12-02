//   Copyright 2019 Freeman Feng<freeman@nuxim.cn>
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package common

import (
	"net"
	"strconv"
	"strings"
)

func GetFreePort() int {
	addr, err := net.ResolveTCPAddr(ProtoTCP, "localhost:0")
	if err != nil {
		return UNDEFINED
	}

	listener, err := net.ListenTCP(ProtoTCP, addr)
	if err != nil {
		return UNDEFINED
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func GetPublicIP() string {
	ip := GetIP(PubIP)
	if ip != EMPTY {
		return ip
	}
	conn, _ := net.Dial(ProtoUDP, GetIP(PubDNS))
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, COLON)
	ip = localAddr[0:idx]
	SetIP(PubIP, ip)
	return ip
}

func IsPortUnavailable(port int) bool {
	listener, err := net.Listen(ProtoTCP, ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	defer listener.Close()
	return false
}
