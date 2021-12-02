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

package parser

import (
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
	"github.com/nuxim/dragon/dragon/parser/common"
	"github.com/nuxim/dragon/dragon/parser/global"
	"github.com/nuxim/dragon/dragon/parser/suites"
	"github.com/nuxim/dragon/dragon/parser/templates"
)

func Run() {
	ch := AnyChannel(ModuleParser)
	cch := make(chan TaskRequest)
	tch := make(chan TaskRequest)
	gch := make(chan TaskRequest)
	rch := make(chan TaskRequest)
	go common.Parse(cch)    // 解析公共变量
	go templates.Parse(tch) // 解析模板
	go global.Parse(gch)    // 解析全局设置
	go suites.Parse(rch)    // 解析用例
	for {
		x := <-ch
		r := x.(TaskRequest)
		cch <- r
		tch <- r
		gch <- r
		rch <- r
		go process(r)
	}
}

func process(r TaskRequest) {
	t := &Testing{Project: r.Project, IsMocking: r.IsMocking, IsHijacking: r.IsHijacking,
		Config:    ConfigType{},
		Templates: map[string]*TemplateType{}, ControlFlow: map[string]*ControlType{},
		Setup: map[string]*ControlType{}, Teardown: map[string]*ControlType{},
		Checks: map[string]*ControlType{}, Suites: map[string]*TestSuite{},
		Funcs: map[string]*PluginType{},
		Bins:  map[string]string{}}
	x := <-AnyChannel(Join(EMPTY, ParserGlobal, ModuleConfig, r.Task))
	if x == nil {
		return
	}
	t.Config = x.(ConfigType)
	loadFuncs(t)
	loadBins(t)
	pc := <-AnyChannel(Join(EMPTY, ParserGlobal, ModuleControlFlow, r.Task))
	if pc == nil {
		return
	}
	t.ControlFlow = pc.(map[string]*ControlType)
	ps := <-AnyChannel(Join(EMPTY, ParserGlobal, ModulePreCondition, r.Task))
	if ps == nil {
		return
	}
	t.Setup = ps.(map[string]*ControlType)
	pt := <-AnyChannel(Join(EMPTY, ParserGlobal, ModulePostCondition, r.Task))
	if pt == nil {
		return
	}
	t.Teardown = pt.(map[string]*ControlType)
	pv := <-AnyChannel(Join(EMPTY, ParserGlobal, ModuleCheck, r.Task))
	if pv == nil {
		return
	}
	t.Checks = pv.(map[string]*ControlType)

	k := <-AnyChannel(Join(EMPTY, ParserTemplates, r.Task))
	if k == nil {
		return
	}
	t.Templates = k.(map[string]*TemplateType)

	tc := <-AnyChannel(Join(EMPTY, ParserTC, r.Task))
	if tc == nil {
		return
	}
	t.Suites = tc.(map[string]*TestSuite)
	filter(t, r)
	showResult(t)
	AnyChannel(Join(EMPTY, ModuleParser, r.Task)) <- t
}

func loadFuncs(t *Testing) {
	for _, p := range t.Config.Plugins {
		for k, _ := range p.Funcs {
			t.Funcs[k] = p
		}
	}
}

func loadBins(t *Testing) {
	for _, p := range t.Config.Plugins {
		for k, v := range p.Bins {
			t.Bins[k] = v
		}
	}
}

func showResult(t *Testing) {
	Log(ModuleParser, "显示解析结果", t)
	for sid, suite := range t.Suites {
		Log(ModuleParser, "++++++用例集", sid, suite)
		for gid, group := range suite.Groups {
			Log(ModuleParser, "++++++用例组", gid, group)
			for tid, tc := range group.Cases {
				Log(ModuleParser, ">>>>> 用例", tid, tc)
				for _, step := range tc.Steps {
					Log(ModuleParser, "++++++++ 用例步骤", step)
					if len(step.Setting.Ops) > 0 {
						Info(ModuleParser, "用例集", sid, "用例组", gid, "用例", tid, "操作", step.Setting.Ops)
					}
					if len(step.Expects) > 0 {
						Info(ModuleParser, "用例集", sid, "用例组", gid, "用例", tid, "期望", step.Expects)
					}
				}
			}
		}
	}
}

func filter(t *Testing, r TaskRequest) {
	filterGroups(t, r)
	filterCases(t, r)
	filterTags(t, r)
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

func filterGroups(t *Testing, r TaskRequest) {
	if len(r.Group) == 0 {
		Log(ModuleParser, "不需要过滤用例组")
		return
	}
	for _, suite := range t.Suites {
		gs := map[string]*GroupType{}
		for gid, g := range suite.Groups {
			for _, key := range r.Group {
				if match(g.Title, key) || match(g.GID, key) {
					gs[gid] = g
				}
			}
		}
		suite.Groups = map[string]*GroupType{}
		if len(gs) > 0 {
			for k, v := range gs {
				suite.Groups[k] = v
			}
		}
	}
}

func filterCases(t *Testing, r TaskRequest) {
	if len(r.Case) == 0 {
		Log(ModuleParser, "不需要过滤用例")

		return
	}
	for _, suite := range t.Suites {
		for _, g := range suite.Groups {
			cs := map[string]*CaseType{}
			for tid, c := range g.Cases {
				for _, key := range r.Case {
					if match(c.Title, key) || match(c.TID, key) {
						cs[tid] = c
					}
				}
			}
			g.Cases = map[string]*CaseType{}
			if len(cs) > 0 {
				for k, v := range cs {
					g.Cases[k] = v
				}
			}
		}
	}
}

func filterTags(t *Testing, r TaskRequest) {
	if len(r.Tags) == 0 {
		Log(ModuleParser, "不需要过滤用例标签")
		return
	}
	for _, suite := range t.Suites {
		for _, g := range suite.Groups {
			cs := map[string]*CaseType{}
			for tid, c := range g.Cases {
				for _, tag := range c.Tags {
					for _, key := range r.Tags {
						if match(tag, key) {
							cs[tid] = c
						}
					}
				}
			}
			g.Cases = map[string]*CaseType{}
			if len(cs) > 0 {
				for k, v := range cs {
					g.Cases[k] = v
				}
			}
		}
	}
}
