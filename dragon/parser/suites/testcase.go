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

	"github.com/FreemanFeng/dragon/dragon/util"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func parseCases(r TaskRequest, start int, content [][]byte, nodes map[string]*GroupType) (int, int, string) {
	var gp *GroupType
	var p *CaseType
	var sp *StepType
	var exp []OpType
	var buf []byte
	var tcTags []string
	var tcPlugins []string
	var tcBugs []string
	var tcFeatures []string
	var groupPlugins []string
	var groupSetups []string
	var groupTeardowns []string
	var groupArgs []OpType
	var groupMatches []OpType
	var tcArgs []OpType
	var tcMatches []OpType
	groupExtends := map[string]string{}
	tcExtends := map[string]string{}
	var tmpArgs []OpType
	var tmpMatches []OpType
	tid := EMPTY
	gid := EMPTY
	groupTitle := EMPTY
	groupDesc := EMPTY
	groupProto := HTTP
	groupTimeout := DefaultTimeout
	tcTitle := EMPTY
	tcDesc := EMPTY
	tcFill := FillRandom
	tcRandom := 0
	tcLocal := 0
	tcRounds := int64(1)
	tcConcurrence := 1
	tcProto := EMPTY
	tcTimeout := DefaultTimeout
	tcAuthor := EMPTY
	tcMaintainer := EMPTY
	tcLevel := EMPTY
	tcVersion := EMPTY
	rw := regexp.MustCompile(`\w+`)
	rs := regexp.MustCompile(`\s+`)
	n := len(content)
	end := start
	code := SUCCESSFUL
	msg := MsgSuccessful
	gp = nil
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
			groupTitle = h[1]
			Debug(ModuleParser, "用例组标题为", groupTitle)
			continue
		}
		if groupTitle != EMPTY && !MatchStart(c, MarkClass) && !MatchStart(c, MarkClass2) {
			groupDesc += string(c)
			b := bytes.TrimSpace(c)
			k := bytes.Split(b, []byte(OpEqual))
			if len(k) > 1 {
				s := string(bytes.TrimSpace(k[0]))
				// 处理参数、扩展属性
				if strings.Index(s, SPACE) > 0 {
					x := strings.Split(s, SPACE)
					s = x[0]
				}
				t := string(bytes.TrimSpace(k[1]))
				h := rs.Split(t, UNDEFINED)
				switch s {
				case ProtoCN, ProtoEN:
					groupProto = h[0]
				case PluginsCN, PluginsEN:
					groupPlugins = h
				case SetupsCN, SetupsEN:
					groupSetups = h
				case TeardownsCN, TeardownsEN:
					groupTeardowns = h
				case TimeoutCN, TimeoutEN:
					groupTimeout = util.CalculateSecs(h[0])
				case ArgsCN, ArgsEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					code, msg, tmpArgs = ParseOp(bytes.TrimSpace(b[x:n]))
					// 合并参数属性值
					groupArgs = append(groupArgs, tmpArgs...)
				case MatchCN, MatchEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					code, msg, tmpMatches = ParseOp(bytes.TrimSpace(b[x:n]))
					// 合并匹配值
					groupMatches = append(groupMatches, tmpMatches...)
				case ExtendCN, ExtendEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					k = bytes.Split(b[x:n], []byte(OpEqual))
					s = string(bytes.TrimSpace(k[0]))
					t = string(bytes.TrimSpace(k[1]))
					groupExtends[s] = t
				}
			}
			continue
		}
		end = i
		// 用例组
		if MatchStart(c, MarkClass) {
			region = UNDEFINED
			exp = nil
			b := rw.Find(c)
			if len(b) == 0 {
				return end, FAILED, MsgIncompleteRule
			}
			gid = string(b)
			_, ok := nodes[gid]
			if ok {
				Debug(ModuleParser, "!!!!!! 用户组ID", gid, "已存在，当前行号", i)
				return end, FAILED, MsgDuplicatedName
			}
			Debug(ModuleParser, "发现用例组", gid)
			nodes[gid] = &GroupType{
				Line:    i + 1,
				Title:   groupTitle,
				Desc:    groupDesc,
				GID:     gid,
				Timeout: groupTimeout,
				Proto:   groupProto,
				PreIDs:  groupSetups,
				PostIDs: groupTeardowns,
				Matches: groupMatches,
				Plugins: groupPlugins,
				Args:    groupArgs,
				Extends: map[string][]string{},
				Cases:   map[string]*CaseType{}}
			gp = nodes[gid]
			for x, y := range groupExtends {
				gp.Extends[x] = rs.Split(y, UNDEFINED)
			}
			groupTitle = EMPTY
			groupDesc = EMPTY
			groupProto = HTTP
			groupTimeout = DefaultTimeout
			groupSetups = nil
			groupTeardowns = nil
			groupArgs = nil
			groupMatches = nil
			groupPlugins = nil
			groupExtends = map[string]string{}
			k := bytes.Split(c, []byte(OpInput))
			b = bytes.TrimSpace(k[1])
			if len(k) > 1 && len(b) > 0 {
				h := bytes.Split(b, []byte(VerticalBar))
				for j, b := range h {
					if j == 0 {
						s := string(bytes.TrimSpace(h[0]))
						gp.CIDs = rs.Split(s, UNDEFINED)
						continue
					}
					k = bytes.Split(b, []byte(OpEqual))
					if len(k) < 2 {
						continue
					}
					s := string(bytes.TrimSpace(k[0]))
					b = bytes.TrimSpace(k[1])
					t := rs.Split(string(b), UNDEFINED)
					switch s {
					case SetupsCN, SetupsEN:
						gp.PreIDs = t
					case TeardownsCN, TeardownsEN:
						gp.PostIDs = t
					case ProtoCN, ProtoEN:
						gp.Proto = t[0]
					case PluginsCN, PluginsEN:
						gp.Plugins = append(gp.Plugins, t...)
					case TimeoutCN, TimeoutEN:
						gp.Timeout = util.CalculateSecs(t[0])
					}
				}
			}
			tid = EMPTY
			Info(ModuleParser, "发现用例组", gid, "标题", gp.Title, "协议",
				gp.Proto, "插件", gp.Plugins, "控制流", gp.CIDs,
				"参数", gp.Args, "前置条件", gp.PreIDs, "后置条件",
				gp.PostIDs, "超时", gp.Timeout, "扩展", gp.Extends)
			continue
		}
		if MatchStart(c, MarkRegion3) {
			region = UNDEFINED
			s := string(bytes.TrimSpace(c))
			h := rs.Split(s, UNDEFINED)
			if len(h) < 2 {
				return end, FAILED, MsgIncompleteRule
			}
			tcTitle = h[1]
			tcProto = gp.Proto
			tcTimeout = gp.Timeout
			Debug(ModuleParser, "发现用例标题", tcTitle)
			continue
		}
		if len(tcTitle) > 0 && !MatchStart(c, MarkClass2) {
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
				case ArgsCN, ArgsEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					code, msg, tmpArgs = ParseOp(bytes.TrimSpace(b[x:n]))
					// 合并参数属性值
					tcArgs = append(tcArgs, tmpArgs...)
				case MatchCN, MatchEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					code, msg, tmpMatches = ParseOp(bytes.TrimSpace(b[x:n]))
					// 合并匹配值
					tcMatches = append(tcMatches, tmpMatches...)
				case TimeoutCN, TimeoutEN:
					tcTimeout = util.CalculateSecs(h[0])
				case TagsCN, TagsEN:
					tcTags = append(tcTags, h...)
				case PluginsCN, PluginsEN:
					tcPlugins = h
				case RoundsCN, RoundsEN:
					x, e := strconv.Atoi(h[0])
					if e == nil {
						tcRounds = int64(x)
					}
				case ConcurrenceCN, ConcurrenceEN:
					x, e := strconv.Atoi(h[0])
					if e == nil {
						tcConcurrence = x
					}
				case FillCN, FillEN:
					switch h[0] {
					case ALL:
						tcFill = FillAll
					case RANDOM:
						tcFill = FillRandom
					}
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
				case AuthorCN, AuthorEN:
					tcAuthor = h[0]
				case MaintainerCN, MaintainerEN:
					tcMaintainer = h[0]
				case BugsCN, BugsEN:
					tcBugs = append(tcBugs, h...)
				case FeaturesCN, FeaturesEN:
					tcFeatures = append(tcFeatures, h...)
				case LevelCN, LevelEN:
					tcLevel = h[0]
				case VersionCN, VersionEN:
					tcVersion = h[0]
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
		// 空行或注释
		if len(bytes.TrimSpace(c)) == 0 || MatchStart(bytes.TrimSpace(c), MarkComment, MarkComment2) {
			continue
		}
		// 判断是否已获取组ID
		_, ok := nodes[gid]
		if gid == EMPTY || len(gp.CIDs) == 0 || !ok {
			return end, FAILED, MsgIncompleteRule
		}
		// 用例开始
		if MatchStart(c, MarkClass2) {
			region = UNDEFINED
			c = bytes.TrimSpace(bytes.TrimPrefix(c, []byte(MarkClass2)))
			k := bytes.Split(c, []byte(SPACE))
			tid = string(k[0])
			gp.Cases[tid] = &CaseType{
				Line:        i + 1,
				Timeout:     tcTimeout,
				Rounds:      tcRounds,
				Concurrence: tcConcurrence,
				Fill:        tcFill,
				Random:      tcRandom,
				Local:       tcLocal,
				TID:         tid,
				Title:       tcTitle,
				Desc:        tcDesc,
				Author:      tcAuthor,
				Maintainer:  tcMaintainer,
				Version:     tcVersion,
				Level:       tcLevel,
				Proto:       tcProto,
				Bugs:        tcBugs,
				Features:    tcFeatures,
				Tags:        tcTags,
				Extends:     map[string][]string{},
				Vars:        map[string]*VarType{},
				Steps:       []StepType{},
			}
			p = gp.Cases[tid]
			p.Args = append(p.Args, gp.Args...)
			p.Args = append(p.Args, tcArgs...)
			p.Plugins = append(p.Plugins, gp.Plugins...)
			p.Plugins = append(p.Plugins, tcPlugins...)
			p.Matches = append(p.Matches, gp.Matches...)
			p.Matches = append(p.Matches, tcMatches...)
			for x, y := range gp.Extends {
				p.Extends[x] = append(p.Extends[x], y...)
			}
			for x, y := range tcExtends {
				p.Extends[x] = rs.Split(y, UNDEFINED)
			}
			tcTitle = EMPTY
			tcDesc = EMPTY
			tcFill = FillRandom
			tcRandom = 0
			tcLocal = 0
			tcRounds = 1
			tcConcurrence = 1
			tcProto = gp.Proto
			tcTimeout = gp.Timeout
			tcAuthor = EMPTY
			tcMaintainer = EMPTY
			tcLevel = EMPTY
			tcVersion = EMPTY
			tcExtends = nil
			tcArgs = nil
			tcPlugins = nil
			tcMatches = nil
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
				case TimeoutCN, TimeoutEN:
					p.Timeout = util.CalculateSecs(h[0])
				case RoundsCN, RoundsEN:
					x, e := strconv.Atoi(h[0])
					if e == nil {
						p.Rounds = int64(x)
					}
				case ConcurrenceCN, ConcurrenceEN:
					x, e := strconv.Atoi(h[0])
					if e == nil {
						p.Concurrence = x
					}
				case FillCN, FillEN:
					switch h[0] {
					case ALL:
						p.Fill = FillAll
					case RANDOM:
						p.Fill = FillRandom
					}
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
			Info(ModuleParser, "发现用例", tid, "标题", p.Title, "作者", p.Author, "维护者", p.Maintainer,
				"版本", p.Version, "等级", p.Level, "协议", p.Proto, "标签", p.Tags, "需求", p.Features,
				"Bugs", p.Bugs, "并发", p.Concurrence, "轮次", p.Rounds, "参数", p.Args, "扩展", p.Extends,
				"用例组", gp.Cases, "插件", p.Plugins)
			continue
		}
		if tid == EMPTY || p == nil {
			return end, FAILED, MsgIncompleteRule
		}
		// 步骤开始
		if MatchStart(c, MarkItem) {
			c = bytes.TrimPrefix(c, []byte(MarkItem))
		}
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
