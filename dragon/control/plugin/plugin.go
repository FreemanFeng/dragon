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

package plugin

import (
	"os/exec"
	"strconv"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func Run() {
	root := GetPath(PluginsPath, DefaultPluginsPath)
	Debug(ModulePlugin, "monitoring plugin folder", root, "...")
	ch := AnyChannel(ModulePlugin)
	for {
		x := <-ch
		s := x.(string)
		go runPlugin(s)
	}
}

func runPlugin(s string) {
	dst := GetDir(s)
	c, files := Decompress(dst, s)
	if c != SUCCESSFUL {
		return
	}
	p := strings.Split(dst, SLASH)
	n := len(p) - 1
	service := Join(DOT, ModulePlugin, p[n])
	name := GetPlugin(service, DefaultPluginName)
	for _, file := range files {
		p = strings.Split(file, SLASH)
		n = len(p) - 1
		if name == p[n] {
			k := Join(SLASH, dst, file)
			go startPlugin(service, k)
			return
		}
	}
	return
}

/*
   refer to https://blog.csdn.net/sinat_36521655/article/details/79296181 for shell script getting param
*/
func startPlugin(service, plugin string) int {
	p := GetPort(service)
	if p == UNDEFINED {
		p = GetFreePort()
	}
	Debug(ModuleControl, "running plugin:", plugin, "-p", p)
	cmd := exec.Command(plugin, "-p", strconv.Itoa(p))
	err := cmd.Start()
	if err != nil {
		Err(ModuleControl, err)
		return FAILED
	}
	return SUCCESSFUL
}
