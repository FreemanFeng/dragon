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

package suites

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
	"github.com/FreemanFeng/dragon/dragon/util"
)

func Parse(ch chan TaskRequest) {
	for {
		r := <-ch
		Debug(ParserTC, "Request", r)
		go start(r)
	}
}

func match(s, key string) bool {
	k := strings.Split(key, OpPlus)
	for _, t := range k {
		if strings.Index(s, t) == UNDEFINED {
			return false
		}
	}
	return true
}

func start(r TaskRequest) {
	path := ParserPath(TPath, r.Project, ParserTC)
	files := util.GetFiles(path, MarkdownType)
	ch := AnyChannel(r.Task)
	i := 0
	code := SUCCESSFUL
	msg := MsgSuccessful
	nodes := map[string]*TestSuite{}
	for _, k := range files {
		s := strings.TrimPrefix(k, path)
		s = strings.TrimPrefix(s, string(os.PathSeparator))
		suite := strings.TrimSuffix(s, MarkdownSuffix)
		ok := false
		for _, key := range r.Suite {
			if match(suite, key) {
				ok = true
			}
		}
		if !ok && len(r.Suite) > 0 {
			continue
		}
		v, err := ioutil.ReadFile(k)
		if err != nil {
			Debug(ParserTC, "Task", r.Task, err)
			continue
		}
		i, code, msg = recognize(r, suite, v, nodes)
		if code != SUCCESSFUL {
			Log(ParserTC, "!!!!!!!!!!!!!!!!!解析失败", msg, "行号", i)
			ch <- TestReport{Stage: PARSING, File: k, Code: code, Line: i + 1, Reason: msg}
			AnyChannel(Join(EMPTY, ParserTC, r.Task)) <- nil
			return
		}
	}
	AnyChannel(Join(EMPTY, ParserTC, r.Task)) <- nodes
	Debug(ParserTC, "path", path, "files", files)
}

func recognize(r TaskRequest, suite string, b []byte, nodes map[string]*TestSuite) (int, int, string) {
	content := bytes.Split(b, []byte(NEWLINE))
	n := len(content)
	code := SUCCESSFUL
	msg := MsgSuccessful
	nodes[suite] = &TestSuite{SID: suite}
	p := nodes[suite]
	Debug(ModuleParser, "recognize", suite)
	for i := 0; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		if !MatchStart(c, MarkRegion) {
			continue
		}
		Debug(ModuleParser, "found", string(c))
		switch {
		// 控制流 Control Flow
		case MatchStart(c, ControlFlow, ControlFlowEN):
			p.ControlFlow = map[string]*ControlType{}
			i, code, msg = ParseFunCall(r, i+1, content, p.ControlFlow, EMPTY)
		// 前置条件 Setup
		case MatchStart(c, SETUP, SetupEN):
			p.Setup = map[string]*ControlType{}
			i, code, msg = ParseFunCall(r, i+1, content, p.Setup, EMPTY)
		// 后置条件 Teardown
		case MatchStart(c, TEARDOWN, TeardownEN):
			p.Teardown = map[string]*ControlType{}
			i, code, msg = ParseFunCall(r, i+1, content, p.Teardown, EMPTY)
		// 校验 Check
		case MatchStart(c, CHECK, CheckEN):
			p.Checks = map[string]*ControlType{}
			i, code, msg = ParseFunCall(r, i+1, content, p.Checks, EMPTY)
		// 流程 Flow
		case MatchStart(c, FLOW, FlowEN):
			p.Flows = map[string]*FlowType{}
			i, code, msg = ParseFlow(r, i+1, content, p.Flows)
		// 用例 Case
		case MatchStart(c, CASE, CaseEN):
			p.Groups = map[string]*GroupType{}
			i, code, msg = parseCases(r, i+1, content, p.Groups)
		default:
			continue
		}
		if code != SUCCESSFUL {
			return i, code, msg
		}
	}
	return n, SUCCESSFUL, MsgSuccessful
}
