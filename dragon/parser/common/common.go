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
	"bytes"
	"io/ioutil"

	"github.com/FreemanFeng/dragon/dragon/parser/suites"

	. "github.com/FreemanFeng/dragon/dragon/common"
	"github.com/FreemanFeng/dragon/dragon/parser/global"
	"github.com/FreemanFeng/dragon/dragon/util"
)

func Parse(ch chan TaskRequest) {
	for {
		r := <-ch
		Debug(ParserGlobal, "Request", r)
		start(r)
	}
}

func start(r TaskRequest) {
	path := ParserPath(TPath, EMPTY, ParserCommon)
	files := util.GetFiles(path, MarkdownType)
	i := 0
	code := SUCCESSFUL
	msg := MsgSuccessful
	config := ConfigType{Actions: map[string]string{}, Attributes: map[string]string{}, Meta: map[string]int{},
		Settings: map[string][]string{}, Nodes: map[string]*NodeType{}, Flows: map[string]*FlowType{},
		Plugins: map[string]*PluginType{}, Messages: map[string]*ConstructType{}, Data: map[string]*ConstructType{}}
	ch := AnyChannel(r.Task)
	for _, k := range files {
		v, err := ioutil.ReadFile(k)
		if err != nil {
			Debug(ParserCommon, "Task", r.Task, err)
			continue
		}
		i, code, msg = recognize(r, v, &config)
		if code != SUCCESSFUL {
			ch <- TestReport{Stage: PARSING, File: k, Code: code, Line: i, Reason: msg}
			AnyChannel(Join(EMPTY, ParserCommon, r.Task)) <- nil
			return
		}
	}
	AnyChannel(Join(EMPTY, ParserCommon, r.Task)) <- config
	Debug(ParserCommon, "path", path, "files", files)
}

func recognize(r TaskRequest, b []byte, p *ConfigType) (int, int, string) {
	content := bytes.Split(b, []byte(NEWLINE))
	n := len(content)
	code := SUCCESSFUL
	msg := MsgSuccessful
	end := UNDEFINED
	for i := 0; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		if !MatchStart(c, MarkRegion) {
			continue
		}
		Debug(ModuleParser, "found", string(c))
		switch {
		// 设置 Setting
		case MatchStart(c, SETTING, SettingEN):
			p.Settings = map[string][]string{}
			i, code, msg = global.ParseRange(r, i+1, content, p.Settings)
		// 节点 Node
		case MatchStart(c, NODE, NodeEN):
			p.Nodes = map[string]*NodeType{}
			i, code, msg = global.ParseNode(r, i+1, content, p.Nodes)
		// 流程 Flow
		case MatchStart(c, FLOW, FlowEN):
			p.Flows = map[string]*FlowType{}
			i, code, msg = suites.ParseFlow(r, i+1, content, p.Flows)
		// 插件 Plugin
		case MatchStart(c, PLUGIN, PluginEN):
			p.Plugins = map[string]*PluginType{}
			i, code, msg = global.ParsePlugin(r, i+1, content, p.Plugins)
		// 消息构造 Message Construct
		case MatchStart(c, MessageConstruct, MessageConstructEN):
			p.Messages = map[string]*ConstructType{}
			i, code, msg = global.ParseMessageConstructor(r, i+1, content, p.Messages)
		// 数据构造 Data Construct
		case MatchStart(c, DataConstruct, DataConstructEN):
			p.Data = map[string]*ConstructType{}
			i, code, msg = global.ParseDataConstructors(r, i+1, content, p.Data)
		}
	}
	return end, code, msg
}
