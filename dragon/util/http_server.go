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

package util

import (
	"crypto/tls"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func RunHttpServer(handler http.HandlerFunc, httpPort, readTimeout, writeTimeout int, quit chan int) {
	Info(ModuleUtil, "Running Http Server on port", httpPort)
	a := []string{COLON, strconv.Itoa(httpPort)}
	h := strings.Join(a, EMPTY)
	s := &http.Server{
		Addr:           h,
		Handler:        handler,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// disable HTTP/2
	s.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	select {
	case <-quit:
		Err(ModuleUtil, "Stop Http Server and Quit")
		return
	}
}
