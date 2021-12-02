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

package src

import (
	"fmt"
	"os"
	"path"
	"plugin"
	"testing"
)

func TestInfo(t *testing.T) {
	//加载插件
	pluginDir := "weather"
	//扫描文件夹下所有so文件
	f, err := os.OpenFile(pluginDir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	fi, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	plugins := make([]os.FileInfo, 0)
	for _, ff := range fi {
		if ff.IsDir() || path.Ext(ff.Name()) != ".so" {
			continue
		}
		plugins = append(plugins, ff)
		pdll, err := plugin.Open(pluginDir + "/" + ff.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		plg, err := pdll.Lookup("Init")
		if err != nil {
			panic(err)
		}
		m := plg.(func() map[string]func(params ...interface{}) interface{})()
		fmt.Println(m)
	}
}
