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
	"plugin"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func compile(r TaskRequest, p *PluginType) {
	var cmd *exec.Cmd
	so := getSo(p)
	if !IsFile(so) || !IsGoCode(p) || !checkNewer(r, p, so) {
		return
	}
	killProcess(r, p)
	Info(ModuleParser, ">>>> 重新编译", so)
	s := []string{"cd ", p.Path, " && ",
		"go build -buildmode=plugin -o ", p.Name, ExtSo, SPACE, p.Name, ExtGo}
	c := strings.Join(s, EMPTY)
	Info(ModuleParser, "执行命令", c)
	cmd = exec.Command("sh", "-c", c)
	out, err := cmd.Output()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		os.Exit(1)
		return
	}
	Info(ModuleParser, string(out))
}

func getSo(p *PluginType) string {
	s := []string{p.Path, string(os.PathSeparator), p.Name, ExtSo}
	return strings.Join(s, EMPTY)
}

func load(r TaskRequest, p *PluginType) {
	//加载插件
	//扫描文件夹下所有so文件
	f, err := os.OpenFile(p.Path, os.O_RDONLY, 0666)
	if err != nil {
		Err(ModuleCommon, "Open so文件出错", err)
		panic(err)
	}
	fi, err := f.Readdir(-1)
	if err != nil {
		Err(ModuleCommon, "读取插件目录出错", err)
		panic(err)
	}

	for _, ff := range fi {
		ext := path.Ext(ff.Name())
		if ff.IsDir() || ext != ExtSo {
			continue
		}
		pdll, err := plugin.Open(path.Join(p.Path, ff.Name()))
		if err != nil {
			if !IsFile(path.Join(p.Path, ff.Name())) {
				Err("！！！！无法识别插件文件")
				os.Exit(1)
			}
			Err(ModuleCommon, "!!!!!无法打开插件", path.Join(p.Path, ff.Name()), err)
			panic(err)
			//continue
		}
		plg, err := pdll.Lookup("Init")
		if err != nil {
			Err(ModuleCommon, "!!!!!定位插件初始化函数失败", err)
			panic(err)
			//return
		}
		if k, ok := plg.(func() map[string]interface{}); ok {
			m := k()
			Info(ModuleCommon, ">>>>>>>> 加载插件完成", m)
			for k, v := range m {
				// 若函数调用列表不为空，将过滤掉没有指定的函数
				if len(p.Calls) > 0 && !HasKey(k, p.Calls) {
					continue
				}
				// 若函数屏蔽列表不为空，将过滤掉命中的函数
				if len(p.Mutes) > 0 && HasKey(k, p.Mutes) {
					continue
				}
				p.Funcs[k] = v
			}
		}
	}
}
