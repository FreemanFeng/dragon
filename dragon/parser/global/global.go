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
		// ?????? Setting
		case MatchStart(c, SETTING, SettingEN):
			i, code, msg = ParseRange(r, i+1, content, p.Settings)
			configSystem(SActions, p.Settings, p.Actions)
			configSystem(SAttrs, p.Settings, p.Attributes)
			Log("????????? Attribute", p.Attributes)
		// ?????? Node
		case MatchStart(c, NODE, NodeEN):
			i, code, msg = ParseNode(r, i+1, content, p.Nodes)
		// ?????? Flow
		case MatchStart(c, FLOW, FlowEN):
			i, code, msg = suites.ParseFlow(r, i+1, content, p.Flows)
		// ?????? Plugin
		case MatchStart(c, PLUGIN, PluginEN):
			i, code, msg = ParsePlugin(r, i+1, content, p.Plugins)
		// ???????????? Message Construct
		case MatchStart(c, MessageConstruct, MessageConstructEN):
			i, code, msg = ParseMessageConstructor(r, i+1, content, p.Messages)
		// ???????????? Data Construct
		case MatchStart(c, DataConstruct, DataConstructEN):
			i, code, msg = ParseDataConstructors(r, i+1, content, p.Data)
		// ????????? Control Flow
		case MatchStart(c, ControlFlow, ControlFlowEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pc, ParserGlobal)
		// ???????????? Setup
		case MatchStart(c, SETUP, SetupEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, ps, ParserGlobal)
		// ???????????? Teardown
		case MatchStart(c, TEARDOWN, TeardownEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pt, ParserGlobal)
		// ?????? Check
		case MatchStart(c, CHECK, CheckEN):
			i, code, msg = suites.ParseFunCall(r, i+1, content, pv, ParserGlobal)
		}
	}
	return end, code, msg
}

func initConfig(p *ConfigType) {
	p.Actions = map[string]string{}
	p.Actions[ActSend] = AKSend        // ????????????
	p.Actions[ActReceive] = AKReceive  // ????????????
	p.Actions[ActCache] = AKCache      // ????????????
	p.Actions[ActDB] = ActKeyDB        // DB??????
	p.Actions[ActHttp] = AKHttp        // HTTP??????
	p.Actions[ActHttps] = AKHttps      // HTTPS??????
	p.Actions[ActMock] = AKMock        // ????????????
	p.Actions[ActAndroid] = AKAndroid  // Android UI?????????
	p.Actions[ActWEB] = ActKeyWEB      // WEB UI?????????
	p.Actions[ActIOS] = ActKeyIOS      // IOS UI?????????
	p.Actions[ActROS] = ActKeyROS      // ROS?????? Robot OS
	p.Actions[ActCOM] = ActKeyCOM      // ????????????
	p.Attributes = map[string]string{} // ????????????
	p.Attributes[Topic] = KeyTopic     // ?????????????????????MQTT
	p.Attributes[Path] = KeyPath       // ???????????????????????????HTTP
	p.Attributes[Args] = KeyArgs       // ???????????????????????????HTTP
	p.Attributes[Header] = KeyHeader   // ????????????????????????HTTP
	p.Attributes[Buffer] = KeyBuffer   // ?????????
	p.Attributes[MainMsg] = KeyMainMsg // ??????????????????????????????????????????msg??????
	p.Attributes[Resp] = KeyResp       // ??????
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
