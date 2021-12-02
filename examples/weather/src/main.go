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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FreemanFeng/dragon/examples/weather/src/control"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], "[-p port][-h][-d]")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var help bool
	var port int
	flag.BoolVar(&help, "h", false, "Show Usage")
	flag.IntVar(&port, "p", 8080, "Control Port")
	flag.Parse()
	if help {
		Usage()
		return
	}
	control.Run(port)
}
