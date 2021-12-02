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
	"regexp"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func ParseNode(r TaskRequest, start int, content [][]byte, nodes map[string]*NodeType) (int, int, string) {
	var sp *NodeType
	var step *StepType
	var buf []byte
	var hosts []string
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
	tunnel := EMPTY
	user := EMPTY
	pass := EMPTY
	path := EMPTY
	method := EMPTY
	proto := EMPTY
	tunuser := EMPTY
	tunpass := EMPTY
	step = nil
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
				Err(ModuleParser, "获取不到服务标题")
				return end, FAILED, MsgIncompleteRule
			}
			title = h[1]
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
				case TunnelCN, TunnelEN:
					tunnel = h[0]
				case TunUserCN, TunUserEN:
					tunuser = h[0]
				case TunPassCN, TunPassEN:
					tunpass = h[0]
				case HostsCN, HostsEN:
					hosts = append(hosts, h...)
				case PluginsCN, PluginsEN:
					plugins = h
				case UserCN, UserEN:
					user = h[0]
				case PassCN, PassEN:
					pass = h[0]
				case PathCN, PathEN:
					path = h[0]
				case MethodCN, MethodEN:
					method = h[0]
				case ProtoCN, ProtoEN:
					proto = h[0]
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
			Debug(ModuleParser, "Find Service ID", id)
			_, ok := nodes[id]
			if ok {
				Err(ModuleParser, "!!!!!!!!!!!服务ID", id, "已存在，当前行号", i)
				return end, FAILED, MsgDuplicatedName
			}
			nodes[id] = &NodeType{NID: id, Title: title, Desc: desc, User: user, Pass: pass,
				Tunnel: tunnel, TunUser: tunuser, TunPass: tunpass, Path: path, Method: method, Proto: proto,
				Extends: map[string][]string{}}
			sp = nodes[id]
			for x, y := range extends {
				sp.Extends[x] = rs.Split(y, UNDEFINED)
			}
			sp.Hosts = append(sp.Hosts, hosts...)
			sp.Plugins = append(sp.Plugins, plugins...)

			title = EMPTY
			desc = EMPTY
			hosts = nil
			plugins = nil
			extends = nil
			k := bytes.Split(c, []byte(VerticalBar))
			for i, v := range k {
				if i == 0 {
					x := bytes.Split(v, []byte(OpInput))
					if len(x) > 1 {
						sp.Args = rs.Split(string(bytes.TrimSpace(x[1])), UNDEFINED)
					}
					continue
				}
				h := bytes.Split(v, []byte(OpEqual))
				if len(h) > 1 {
					s := string(bytes.TrimSpace(h[0]))
					t := string(bytes.TrimSpace(h[1]))
					a := rs.Split(t, UNDEFINED)
					switch s {
					case HostsCN, HostsEN:
						sp.Hosts = append(sp.Hosts, a...)
					case PathCN, PathEN:
						sp.Path = a[0]
					case MethodCN, MethodEN:
						sp.Method = a[0]
					case ProtoCN, ProtoEN:
						sp.Proto = a[0]
					case PluginsCN, PluginsEN:
						sp.Plugins = append(sp.Plugins, a...)
					}
				}
			}
			Info(ModuleParser, "发现控制流", id, "标题", sp.Title, "主机", sp.Hosts, "用户", sp.User, "插件", sp.Plugins,
				"密码", sp.Pass, "通道", sp.Tunnel, "通道用户", sp.TunUser, "通道密码", sp.TunPass, "扩展", sp.Extends)
			continue
		}
		// 空行或注释
		if len(bytes.TrimSpace(c)) == 0 || MatchStart(bytes.TrimSpace(c), MarkComment, MarkComment2) {
			continue
		}
		// 判断是否已获取控制流ID
		_, ok := nodes[id]
		if id == EMPTY || !ok || sp == nil {
			return end, FAILED, MsgIncompleteRule
		}
		// 控制流开始
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
		t := len(sp.Steps)
		sp.Steps = append(sp.Steps, StepType{Line: i + 1})
		step = &sp.Steps[t]

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
		code, msg, step.Setting = ParseSetting(buf)
	}

	return end, SUCCESSFUL, MsgSuccessful
}
