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
	"regexp"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func ParseFunCall(r TaskRequest, start int, content [][]byte, nodes map[string]*ControlType, upper string) (int, int, string) {
	var cp *ControlType
	var sp *StepType
	var buf []byte
	var plugins []string
	rw := regexp.MustCompile(`\w+`)
	rs := regexp.MustCompile(`\s+`)
	n := len(content)
	code := SUCCESSFUL
	msg := MsgSuccessful
	end := start
	id := EMPTY
	title := EMPTY
	desc := EMPTY
	// 默认范围
	scope := ScopeSuite
	cp = nil
	sp = nil
	extends := map[string]string{}
	for i := start; i < n; i++ {
		if code != SUCCESSFUL {
			return end, code, msg
		}
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		if MatchStart(c, MarkRegion2) {
			s := string(bytes.TrimSpace(c))
			h := rs.Split(s, UNDEFINED)
			if len(h) < 2 {
				return end, FAILED, MsgIncompleteRule
			}
			title = h[1]
			Debug(ModuleParser, "发现控制流标题", title)
			continue
		}
		if title != EMPTY && !MatchStart(c, MarkClass) {
			desc += string(c)
			b := bytes.TrimSpace(c)
			k := bytes.Split(b, []byte(OpEqual))
			if len(k) > 1 {
				s := string(bytes.TrimSpace(k[0]))
				// 处理扩展属性
				if strings.Index(s, SPACE) > 0 {
					x := strings.Split(s, SPACE)
					s = x[0]
				}
				t := string(bytes.TrimSpace(k[1]))
				h := rs.Split(t, UNDEFINED)
				switch s {
				case PluginsCN, PluginsEN:
					plugins = h
				case ScopeCN, ScopeEN:
					scope = h[0]
				case ExtendCN, ExtendEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					k = bytes.Split(b[x:n], []byte(OpEqual))
					s = string(bytes.TrimSpace(k[0]))
					t = string(bytes.TrimSpace(k[1]))
					extends[s] = t
				}
			}

			continue
		}
		end = i
		if MatchStart(c, MarkClass) {
			id = string(rw.Find(c))
			Debug(ModuleParser, "Find CF ID", id)
			_, ok := nodes[id]
			if ok {
				Debug(ModuleParser, "!!!!!!!!!!!控制流ID", id, "已存在，当前行号", i)
				return end, FAILED, MsgDuplicatedName
			}
			nodes[id] = &ControlType{
				CID:     id,
				Title:   title,
				Desc:    desc,
				Plugins: plugins,
				Scope:   scope,
				Args:    map[string]int{},
				Extends: map[string][]string{},
				Vars:    map[string]*VarType{}}
			cp = nodes[id]
			for x, y := range extends {
				cp.Extends[x] = rs.Split(y, UNDEFINED)
			}
			// 在global目录下定义的控制流，都应该作为全局定义
			if upper == ParserGlobal {
				cp.Scope = ScopeAll
			}
			title = EMPTY
			desc = EMPTY
			scope = EMPTY
			plugins = nil
			k := bytes.Split(c, []byte(VerticalBar))
			for i, v := range k {
				if i == 0 {
					x := bytes.Split(v, []byte(OpInput))
					if len(x) > 1 {
						s := rs.Split(string(bytes.TrimSpace(x[1])), UNDEFINED)
						for _, h := range s {
							cp.Args[h] = 1
						}
					}
					continue
				}
				h := bytes.Split(v, []byte(OpEqual))
				if len(h) > 1 {
					s := string(bytes.TrimSpace(h[0]))
					t := string(bytes.TrimSpace(h[1]))
					a := rs.Split(t, UNDEFINED)
					switch s {
					case PluginsCN, PluginsEN:
						cp.Plugins = append(cp.Plugins, a...)
					case ScopeCN, ScopeEN:
						cp.Scope = a[0]
					}
				}
			}
			Info(ModuleParser, "发现控制流", cp.CID, "标题", cp.Title, "插件", cp.Plugins,
				"范围", cp.Scope, "参数", cp.Args, "扩展", cp.Extends)
			continue
		}
		// 空行或注释
		if len(bytes.TrimSpace(c)) == 0 || MatchStart(bytes.TrimSpace(c), MarkComment, MarkComment2) {
			continue
		}
		// 判断是否已获取控制流ID
		_, ok := nodes[id]
		if id == EMPTY || !ok || cp == nil {
			return end, FAILED, MsgIncompleteRule
		}
		// 控制流开始,*可有可无
		if MatchStart(c, MarkItem) {
			c = bytes.TrimPrefix(c, []byte(MarkItem))
		}
		// 空行或注释
		if len(bytes.TrimSpace(c)) == 0 || MatchStart(bytes.TrimSpace(c), MarkComment2) {
			continue
		}
		if MatchInside(c, MarkComment2) {
			x := Find(c, MarkComment2)
			c = bytes.TrimSpace(c[:x])
		}
		t := len(cp.Flows)
		cp.Flows = append(cp.Flows, StepType{Line: i + 1})
		sp = &cp.Flows[t]

		// 设置
		buf = bytes.TrimSpace(c)
		if MatchInside(c, NotPresent, ExtraPresent, OpConstruct) {
			for k := i + 1; k < n; k++ {
				c := bytes.TrimSpace(content[k])
				if !MatchInside(c, MarkRegion, MarkClass, MarkItem, OpEqual, NotEqual, MoreThan, LessThan,
					OpWithout, OpConstruct, NotPresent, ExtraPresent) {
					h := [][]byte{buf, c}
					buf = bytes.Join(h, []byte(SPACE))
					i = k
					continue
				}
				break
			}
		}
		code, msg, sp.Setting = ParseSetting(buf)
	}

	return end, SUCCESSFUL, MsgSuccessful
}
