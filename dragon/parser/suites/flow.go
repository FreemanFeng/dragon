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
	"strconv"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func ParseFlow(r TaskRequest, start int, content [][]byte, nodes map[string]*FlowType) (int, int, string) {
	var p *FlowType
	var sp *StepType
	var exp []OpType
	var buf []byte
	var tcTags []string
	var tcPlugins []string
	tcExtends := map[string]string{}
	fid := EMPTY
	tcTitle := EMPTY
	tcRandom := 0
	tcLocal := 0
	tcDesc := EMPTY
	tcProto := EMPTY
	rs := regexp.MustCompile(`\s+`)
	n := len(content)
	end := start
	code := SUCCESSFUL
	msg := MsgSuccessful
	p = nil
	sp = nil
	exp = nil
	region := UNDEFINED
	for i := start; i < n; i++ {
		if code != SUCCESSFUL {
			Debug(ModuleParser, "!!!!!!!!! 解析失败", msg, end)
			return end, code, msg
		}
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		if MatchStart(c, MarkRegion2) {
			region = UNDEFINED
			s := string(bytes.TrimSpace(c))
			h := rs.Split(s, UNDEFINED)
			if len(h) < 2 {
				return end, FAILED, MsgIncompleteRule
			}
			tcTitle = h[1]
			Debug(ModuleParser, "流程标题为", tcTitle)
			continue
		}
		if len(tcTitle) > 0 && !MatchStart(c, MarkClass) {
			tcDesc += string(c)
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
				case ProtoCN, ProtoEN:
					tcProto = h[0]
				case RandomCN, RandomEN:
					x, e := strconv.Atoi(h[0])
					if e == nil && x > 0 {
						tcRandom = x
					}
				case LocalCN, LocalEN:
					x, e := strconv.Atoi(h[0])
					if e == nil && x == 1 {
						tcLocal = x
					}
				case TagsCN, TagsEN:
					tcTags = append(tcTags, h...)
				case PluginsCN, PluginsEN:
					tcPlugins = h
				case ExtendCN, ExtendEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					k = bytes.Split(b[x:n], []byte(OpEqual))
					s = string(bytes.TrimSpace(k[0]))
					t = string(bytes.TrimSpace(k[1]))
					tcExtends[s] = t
				}
			}
			continue
		}

		end = i
		// 空行或注释
		if len(bytes.TrimSpace(c)) == 0 || MatchStart(bytes.TrimSpace(c), MarkComment, MarkComment2) {
			continue
		}
		// 流程开始
		if MatchStart(c, MarkClass) {
			region = UNDEFINED
			c = bytes.TrimSpace(bytes.TrimPrefix(c, []byte(MarkClass)))
			k := bytes.Split(c, []byte(SPACE))
			fid = string(k[0])
			nodes[fid] = &FlowType{
				Line:    i + 1,
				Random:  tcRandom,
				Local:   tcLocal,
				FID:     fid,
				Title:   tcTitle,
				Desc:    tcDesc,
				Proto:   tcProto,
				Tags:    tcTags,
				Extends: map[string][]string{},
				Steps:   []StepType{},
				Vars:    map[string]*VarType{},
			}
			p = nodes[fid]
			p.Plugins = append(p.Plugins, tcPlugins...)
			for x, y := range tcExtends {
				p.Extends[x] = rs.Split(y, UNDEFINED)
			}
			tcTitle = EMPTY
			tcDesc = EMPTY
			tcProto = EMPTY
			tcRandom = 0
			tcLocal = 0
			tcExtends = nil
			tcPlugins = nil
			c = bytes.TrimSpace(bytes.TrimPrefix(c, k[0]))
			k = bytes.Split(c, []byte(VerticalBar))
			for _, b := range k {
				t := bytes.Split(b, []byte(OpEqual))
				if len(t) < 2 {
					continue
				}
				s := string(bytes.TrimSpace(t[1]))
				h := rs.Split(s, UNDEFINED)
				s = string(bytes.TrimSpace(t[0]))
				switch s {
				case TagsCN, TagsEN:
					p.Tags = append(p.Tags, h...)
				case PluginsCN, PluginsEN:
					p.Plugins = append(p.Plugins, h...)
				case ProtoCN, ProtoEN:
					// 有可能有用户自定义协议
					p.Proto = h[0]
				case RandomCN, RandomEN:
					x, e := strconv.Atoi(h[0])
					if e == nil && x > 0 {
						p.Random = x
					}
				case LocalCN, LocalEN:
					x, e := strconv.Atoi(h[0])
					if e == nil && x == 1 {
						p.Local = x
					}
				}

			}
			Info(ModuleParser, "发现流程", fid, "标题", p.Title, "协议", p.Proto, "标签", p.Tags, "扩展", p.Extends,
				"插件", p.Plugins)
			continue
		}
		if fid == EMPTY || p == nil {
			return end, FAILED, MsgIncompleteRule
		}
		// 步骤开始
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
		top := len(p.Steps)
		p.Steps = append(p.Steps, StepType{Line: i + 1})
		// 预期行为
		if MatchInside(c, MarkExpect) || region != UNDEFINED {
			if region == UNDEFINED {
				region = top
				code, msg, c, exp = ParseExpect(c)
			} else {
				code, msg, exp = ParseOp(c)
			}
			for i := range exp {
				TrimOp(&exp[i])
			}
			p.Steps[top].IsExpect = true
			p.Steps[top].Expects = append(p.Steps[top].Expects, exp...)

			continue
		}
		sp = &p.Steps[top]
		// 设置
		buf = bytes.TrimSpace(c)
		if MatchInside(c, NotPresent, ExtraPresent, OpConstruct) {
			for k := i + 1; k < n; k++ {
				b := bytes.TrimSpace(content[k])
				if string(b) == EMPTY || MatchStart(b, MarkComment, MarkComment2) {
					continue
				}
				if !MatchInside(b, MarkRegion, MarkRegion2, MarkRegion3, MarkClass, MarkItem,
					OpEqual, NotEqual, MoreThan, LessThan, SpaceOn, SpaceWith, SpaceWithout, OpConstruct,
					NotPresent, ExtraPresent, SpaceOR) {
					h := [][]byte{buf, b}
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
