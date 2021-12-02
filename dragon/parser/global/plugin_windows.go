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
package global

import (
	"os"
	"os/exec"
	"path"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func compileDLL(r TaskRequest, p *PluginType) bool {
	var cmd *exec.Cmd
	s := []string{p.Name, ".go"}
	k := strings.Join(s, EMPTY)
	s = []string{p.Path, k}
	h := path.Join(s...)
	if !IsFile(h) {
		Err(ModuleParser, "不存在", h)
		os.Exit(1)
		return false
	}
	s = []string{"go build -buildmode=c-shared -o ",
		p.Path, string(os.PathSeparator), p.Name, ".dll ",
		p.Path, string(os.PathSeparator), p.Name, ".go"}
	c := strings.Join(s, EMPTY)
	Info(ModuleParser, "执行命令", c)
	cmd = exec.Command("cmd", "/C", c)
	out, err := cmd.Output()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		os.Exit(1)
		return false
	}
	Info(ModuleParser, string(out))
	return true
}
