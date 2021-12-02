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
	"bytes"
	"io/ioutil"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
	"github.com/FreemanFeng/dragon/dragon/parser/suites"
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
	path := ParserPath(TPath, r.Project, ParserGlobal)
	files := util.GetFiles(path, MarkdownType)
	i := 0
	code := SUCCESSFUL
	msg := MsgSuccessful
	x := <-AnyChannel(Join(EMPTY, ParserCommon, r.Task))
	config := x.(ConfigType)
	ch := AnyChannel(r.Task)
	pc := map[string]*ControlType{}
	ps := map[string]*ControlType{}
	pt := map[string]*ControlType{}
	pv := map[string]*ControlType{}
	initConfig(&config)
	for _, k := range files {
		v, err := ioutil.ReadFile(k)
		if err != nil {
			Debug(ParserGlobal, "Task", r.Task, err)
			continue
		}
		i, code, msg = recognize(r, v, &config, pc, ps, pt, pv)
		if code != SUCCESSFUL {
			Err(msg)
			ch <- TestReport{Stage: PARSING, File: k, Code: code, Line: i, Reason: msg}
			AnyChannel(Join(EMPTY, ParserGlobal, r.Task)) <- nil
			return
		}
	}
	AnyChannel(Join(EMPTY, ParserGlobal, ModuleConfig, r.Task)) <- config
	AnyChannel(Join(EMPTY, ParserGlobal, ModuleControlFlow, r.Task)) <- pc
	AnyChannel(Join(EMPTY, ParserGlobal, ModulePreCondition, r.Task)) <- ps
	AnyChannel(Join(EMPTY, ParserGlobal, ModulePostCondition, r.Task)) <- pt
	AnyChannel(Join(EMPTY, ParserGlobal, ModuleCheck, r.Task)) <- pv
	Debug(ParserGlobal, "path", path, "files", files)
}

func recognize(r TaskRequest, b []byte, p *ConfigType, pc, ps, pt, pv map[string]*ControlType) (int, int, string) {
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
		Log(ModuleParser, "found", string(c))
		switch {
		// 设置 Setting
		case MatchStart(c, SETTING, SettingEN):
			i, code, msg = ParseRange(r, i+1, content, p.Settings)
			configSystem(SActions, p.Settings, p.Actions)
			configSystem(SAttrs, p.Settings, p.Attributes)
			Log("修改后 Attribute", p.Attributes)
		// 节点 Node
		case MatchStart(c, NODE, NodeEN):
			i, code, msg = ParseNode(r, i+1, content, p.Nodes)
		// 流程 Flow
		case MatchStart(c, FLOW, FlowEN):
			i, code, msg = suites.ParseFlow(r, i+1, content, p.Flows)
		// 插件 Plugin
		case MatchStart(c, PLUGIN, PluginEN):
			i, code, msg = ParsePlugin(r, i+1, content, p.Plugins)
		// 消息构造 Message Construct
		case MatchStart(c, MessageConstruct, MessageConstructEN):
			i, code, msg = ParseMessageConstructor(r, i+1, content, p.Messages)
		// 数据构造 Data Construct
		case MatchStart(c, DataConstruct, DataConstructEN):
			i, code, msg = ParseDataConstructors(r, i+1, content, p.Data)
		// 控制流 Control Flow
		case MatchStart(c, ControlFlow, ControlFlowEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pc, ParserGlobal)
		// 前置条件 Setup
		case MatchStart(c, SETUP, SetupEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, ps, ParserGlobal)
		// 后置条件 Teardown
		case MatchStart(c, TEARDOWN, TeardownEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pt, ParserGlobal)
		// 校验 Check
		case MatchStart(c, CHECK, CheckEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pv, ParserGlobal)
		}
	}
	return end, code, msg
}

func initConfig(p *ConfigType) {
	p.Actions = map[string]string{}
	p.Actions[ActSend] = AKSend        // 发送消息
	p.Actions[ActReceive] = AKReceive  // 接收消息
	p.Actions[ActCache] = AKCache      // 缓存操作
	p.Actions[ActDB] = ActKeyDB        // DB操作
	p.Actions[ActHttp] = AKHttp        // HTTP请求
	p.Actions[ActHttps] = AKHttps      // HTTPS请求
	p.Actions[ActMock] = AKMock        // 模拟服务
	p.Actions[ActAndroid] = AKAndroid  // Android UI自动化
	p.Actions[ActWEB] = ActKeyWEB      // WEB UI自动化
	p.Actions[ActIOS] = ActKeyIOS      // IOS UI自动化
	p.Actions[ActROS] = ActKeyROS      // ROS消息 Robot OS
	p.Actions[ActCOM] = ActKeyCOM      // 串口消息
	p.Attributes = map[string]string{} // 属性配置
	p.Attributes[Topic] = KeyTopic     // 主题，一般用于MQTT
	p.Attributes[Path] = KeyPath       // 请求路径，一般用于HTTP
	p.Attributes[Args] = KeyArgs       // 请求参数，一般用于HTTP
	p.Attributes[Header] = KeyHeader   // 请求头，一般用于HTTP
	p.Attributes[Buffer] = KeyBuffer   // 请求体
	p.Attributes[MainMsg] = KeyMainMsg // 消息体，请求的主体内容，对应msg目录
	p.Attributes[Resp] = KeyResp       // 应答
}

func configSystem(key string, m map[string][]string, p map[string]string) {
	k, ok := m[key]
	if ok {
		for _, v := range k {
			h := strings.Split(v, COLON)
			if len(h) < 2 {
				continue
			}
			id := strings.TrimSpace(h[0])
			s := strings.TrimSpace(h[1])
			p[s] = id
		}
	}
	delete(m, key)
}
