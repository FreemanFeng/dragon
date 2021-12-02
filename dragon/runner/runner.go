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

package runner

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nuxim/dragon/dragon/rule"

	"github.com/nuxim/dragon/dragon/util"

	. "github.com/nuxim/dragon/dragon/common"
)

func Run() {
	ch := AnyChannel(ModuleRunner)
	for {
		x := <-ch
		k := x.(TaskRequest)
		go process(k)
		Debug(ModuleRunner, k.Task, k.Project, k.Suite, k.Group, k.Case)
	}
}

func process(r TaskRequest) {
	var gs []*GroupData
	var gd *GroupData
	rand.Seed(time.Now().Unix())
	x := <-AnyChannel(Join(EMPTY, ModuleBuilder, r.Task))
	t := x.(*Testing)
	Debug(ModuleRunner, t)
	ch := AnyChannel(Join(EMPTY, ModuleRunner, r.Task))
	m := map[string]chan *GroupData{}
	done := make(chan int)
	total := 0
	begin := util.Seconds()
	for {
		select {
		case <-done:
			total -= 1
			if total <= 0 {
				Log(ModuleRunner, "执行完毕")
				if r.Duration > 0 && r.Duration < util.Seconds()-begin || r.Duration == 0 {
					Log(ModuleRunner, "执行全局Teardown")
					teardown(r, t, gd, ScopeAll)
					return
				} else if r.Duration > 0 {
					total += len(gs)
					go func() {
						for _, g := range gs {
							m[g.Host] <- g
						}
					}()
				}
			}
		case x := <-ch:
			g := x.(*GroupData)
			_, ok := m[g.Host]
			if !ok {
				m[g.Host] = make(chan *GroupData, MaxQueue)
				go dispatch(r, t, m[g.Host], done)
			}
			total += 1
			if g.Scope == ScopeAll {
				gd = g
				if SUCCESSFUL != setup(r, t, gd, ScopeAll) {
					Err(ModuleRunner, "执行全局setup失败")
					return
				}
				continue
			}
			Log(ModuleRunner, "执行用例组数据", g)
			m[g.Host] <- g
			gs = append(gs, g)
		}
	}

}

func setup(r TaskRequest, t *Testing, gd *GroupData, scope string) int {
	// 模拟或劫持服务无需前置条件
	if r.IsMocking || r.IsHijacking {
		return SUCCESSFUL
	}
	tch := AnyChannel(r.Task)
	// 执行全局setup
	if scope == ScopeAll {
		for _, v := range t.Setup {
			for _, flow := range v.Flows {
				if SUCCESSFUL != runStep(UNDEFINED, r, t, nil, nil, v, &flow.Setting, nil, gd.CFs[v.CID], nil) {
					Err(ModuleRunner, "全局Setup执行不成功")
					tch <- TestReport{Stage: EXECUTING, Reason: MsgSetupFailed}
					return FAILED
				}
			}
		}
		return SUCCESSFUL
	}
	// 执行用例集或用例组或用例的setup
	suite, _ := t.Suites[gd.SID]
	group, _ := suite.Groups[gd.GID]
	for _, sid := range group.PreIDs {
		v, ok := suite.Setup[sid]
		if !ok || scope != v.Scope {
			continue
		}
		for _, flow := range v.Flows {
			if SUCCESSFUL != runStep(UNDEFINED, r, t, nil, nil, v, &flow.Setting, nil, gd.CFs[v.CID], nil) {
				Err(ModuleRunner, "用例组、用例集、用例Setup执行不成功")
				tch <- TestReport{Stage: EXECUTING, Reason: MsgSetupFailed}
				return FAILED
			}
		}
	}
	return SUCCESSFUL
}

func teardown(r TaskRequest, t *Testing, gd *GroupData, scope string) int {
	// 模拟或劫持服务无需后置条件
	if r.IsMocking || r.IsHijacking {
		return SUCCESSFUL
	}
	tch := AnyChannel(r.Task)
	// 执行全局的teardown
	if scope == ScopeAll {
		for _, v := range t.Teardown {
			for _, flow := range v.Flows {
				if SUCCESSFUL != runStep(UNDEFINED, r, t, nil, nil, v, &flow.Setting, nil, gd.CFs[v.CID], nil) {
					Log(ModuleRunner, "执行全局Teardown失败")
					tch <- TestReport{Stage: EXECUTING, Reason: MsgTeardownFailed}
					return FAILED
				}
			}
		}
		return SUCCESSFUL
	}
	// 执行用例集或用例组或用例的teardown
	suite, _ := t.Suites[gd.SID]
	group, _ := suite.Groups[gd.GID]
	for _, sid := range group.PostIDs {
		v, ok := suite.Teardown[sid]
		if !ok || scope != v.Scope {
			continue
		}
		for _, flow := range v.Flows {
			if SUCCESSFUL != runStep(UNDEFINED, r, t, nil, nil, v, &flow.Setting, nil, gd.CFs[v.CID], nil) {
				Log(ModuleRunner, "执行用例Teardown失败")
				tch <- TestReport{Stage: EXECUTING, Reason: MsgTeardownFailed}
				return FAILED
			}
		}
	}
	return SUCCESSFUL
}

func dispatch(r TaskRequest, t *Testing, ch chan *GroupData, done chan int) {
	for {
		gd := <-ch
		start(r, t, gd)
		done <- 1
	}
}

func start(r TaskRequest, t *Testing, gd *GroupData) {
	Debug(ModuleRunner, "Received Group Data", gd.Host, gd.SID, gd.GID)
	s, _ := t.Suites[gd.SID]
	g, _ := s.Groups[gd.GID]
	if SUCCESSFUL != setup(r, t, gd, ScopeSuite) {
		Err(ModuleRunner, "执行用例集Setup不成功")
		return
	}
	if SUCCESSFUL != setup(r, t, gd, ScopeGroup) {
		Err(ModuleRunner, "执行用例组Setup不成功")
		return
	}
	for k, c := range g.Cases {
		cd, ok := gd.Cases[k]
		if !ok {
			continue
		}
		ch := make(chan int)
		if SUCCESSFUL != setup(r, t, gd, ScopeCase) {
			Err(ModuleRunner, "执行用例Setup不成功")
			continue
		}
		for i := 0; i < c.Concurrence; i++ {
			// 并发id从1开始, 控制流和流程数据需要复制
			go run(i+1, r, t, s, g, c, gd, cd[i], ch)
		}
		// 等待所有请求完成
		for i := 0; i < c.Concurrence; i++ {
			<-ch
		}
		Log(ModuleRunner, "跑完用例，执行用例Teardown")
		teardown(r, t, gd, ScopeCase)
	}
	teardown(r, t, gd, ScopeGroup)
	teardown(r, t, gd, ScopeSuite)
}

func run(id int, r TaskRequest, t *Testing, s *TestSuite, g *GroupType, c *CaseType,
	gd *GroupData, cd *VarData, ch chan int) {
	begin := util.Seconds()
	endless := false
	if c.Rounds == 0 {
		endless = true
	}
	Log(ModuleRunner, "总的执行轮次", c.Rounds)
	total := int64(0)
	if endless && c.Timeout == 0 || !endless && c.Rounds <= 0 {
		Log(ModuleRunner, "无法执行用例，请设置执行时长或执行次数!")
		return
	}
	rounds := c.Rounds
	for {
		if endless && c.Timeout > 0 && util.Seconds()-begin > c.Timeout || rounds <= 0 && !endless {
			Log(ModuleRunner, "是否无限执行", endless, "超时时间", c.Timeout, "当前剩余执行轮次", rounds)
			break
		}

		if SUCCESSFUL != runOnce(id, r, t, s, g, c, gd, cd) {
			Err(ModuleRunner, "用例集", s.SID, "用例组", g.GID, "用例", c.TID, "执行失败!")
		} else {
			Info(ModuleRunner, "用例集", s.SID, "用例组", g.GID, "用例", c.TID, "执行成功")
		}

		rounds -= 1
		total += 1
		Log(ModuleRunner, "用例执行次数", total)
	}
	ch <- DONE
}

func runOnce(id int, r TaskRequest, t *Testing, s *TestSuite, g *GroupType, c *CaseType, gd *GroupData, cd *VarData) int {
	cfs := map[string]*VarData{}
	fds := map[string]*VarData{}
	CopyVarDataMap(gd.CFs, cfs)
	CopyVarDataMap(gd.FDs, fds)
	for _, cid := range g.CIDs {
		pc, ok := t.ControlFlow[cid]
		if !ok {
			pc, ok = s.ControlFlow[cid]
		}
		if !ok {
			Log(ModuleRunner, "无法找到控制流ID", cid)
			return FAILED
		}
		if SUCCESSFUL != runControlFlow(id, r, t, s, c, pc, cd, cfs[cid], fds) {
			return FAILED
		}
	}
	return SUCCESSFUL
}

// 执行控制流
func runControlFlow(id int, r TaskRequest, t *Testing, s *TestSuite, c *CaseType, pc *ControlType, cd, pd *VarData, fds map[string]*VarData) int {
	for _, k := range pc.Flows {
		if SUCCESSFUL != runStep(id, r, t, s, c, pc, &k.Setting, cd, pd, fds) {
			return FAILED
		}
	}
	// 只有用例里头才有校验逻辑
	for _, k := range c.Steps {
		if SUCCESSFUL != runVerification(id, r, t, c, pc, &k, cd, pd, fds) {
			return FAILED
		}
	}
	return SUCCESSFUL
}

// 执行步骤
func runStep(id int, r TaskRequest, t *Testing, st *TestSuite, ct *CaseType, pc *ControlType,
	sp *SettingType, cd, pd *VarData, fds map[string]*VarData) int {
	if len(sp.Ops) == 0 {
		return SUCCESSFUL
	}
	p := &sp.Ops[0]
	if !rule.FoundID(id, p.CIDs) {
		return SUCCESSFUL
	}
	Log(ModuleRunner, "执行步骤", sp.Ops)
	if !rule.ConditionSatisfied(sp.Ops, pc.Vars, pd) {
		return SUCCESSFUL
	}
	for i := range sp.Ops {
		p = &sp.Ops[i]
		// 条件前的操作已完成
		if p.Op == OpOn {
			break
		}
		// 调用函数
		if len(p.Value) > 0 && !rule.IsConfig(p.Value[0], t) && !IsDigit(p.Value[0]) && rule.IsCall(p.Value[0], t) {
			key := p.Value[0]
			if SUCCESSFUL != runReserved(id, key, r, t, st, ct, pc, p, cd, pd, fds) {
				return FAILED
			}
			rule.RunBuiltin(CFVar, EMPTY, EMPTY, key, r, t, ct, pc, p, pd)
			rule.RunFuncs(CFVar, EMPTY, EMPTY, key, r, t, ct, pc, p, pd)
			if SUCCESSFUL != rule.RunBins(CFVar, EMPTY, EMPTY, key, r, t, ct, pc, p, pd) {
				return FAILED
			}
			continue
		}
		rule.SetValue(EMPTY, EMPTY, ct, pd, pc.Vars, p)
	}
	return SUCCESSFUL
}

func runVerification(id int, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType, sp *StepType, cd, pd *VarData, fds map[string]*VarData) int {
	if len(sp.Expects) == 0 {
		return SUCCESSFUL
	}
	p := &sp.Expects[0]
	if !rule.FoundID(id, p.CIDs) {
		return SUCCESSFUL
	}
	Log(ModuleRunner, "执行校验", sp.Expects)
	k := rule.SearchOp(OpOn, sp.Expects)
	if k == UNDEFINED {
		k = len(sp.Expects)
	}
	if !rule.VerificationPassed(k, sp.Expects, ct.Vars, cd, fds) {
		return FAILED
	}
	return SUCCESSFUL
}

func runReserved(id int, key string, r TaskRequest, t *Testing, st *TestSuite,
	ct *CaseType, pc *ControlType, op *OpType, cd, pd *VarData, fds map[string]*VarData) int {
	if !IsReserved(key) {
		return SUCCESSFUL
	}
	Log(ModuleRunner, "++++++++++++++++ 调用保留函数 +++++++++++++++++++")
	// 遍历所有消息
	for pos, v := range op.Value {
		if pos == 0 {
			continue
		}
		p, ok := pc.Vars[v]
		vt := CFVar
		if !ok {
			if ct == nil {
				Err(ModuleRunner, "无法识别变量名", v)
				return FAILED
			}
			if isArg(v, pc.Args) {
				p, ok = ct.Vars[v]
				if !ok {
					Err(ModuleRunner, "无法识别变量名", v)
					return FAILED
				}
				vt = CaseVar
			} else {
				Err(ModuleRunner, "变量", v, "无法识别!")
				return FAILED
			}
		}
		Log(ModuleRunner, "变量", v, p, "Key", p.Key, "节点", p.Node, "模板数", len(p.Templates))
		if SUCCESSFUL != runCaseTemplates(id, vt, key, r, t, st, ct, pc, op, cd, pd, fds, p) {
			return FAILED
		}
	}
	return SUCCESSFUL
}

func runCaseTemplates(id, vt int, key string, r TaskRequest, t *Testing, st *TestSuite,
	ct *CaseType, pc *ControlType, op *OpType, cd, pd *VarData, fds map[string]*VarData, p *VarType) int {
	var vs []*TemplateVarType
	var fd *VarData
	var ft *FlowType
	n := len(p.Templates)
	if ct.Random > n {
		n = ct.Random
	}
	for n > 0 {
		// 取指针
		vs = nil
		a := rand.Intn(len(p.Templates))
		if ct.Random == 0 {
			a = len(p.Templates) - n
		}
		tv := &p.Templates[a]
		n--
		total := 1
		// 消息模板组或流程
		if tv.IsVar && !tv.IsFlow {
			total = tv.Count
			x, ok := pc.Vars[tv.Key]
			if !ok || vt == CaseVar {
				x, ok = ct.Vars[tv.Key]
				if !ok {
					Err(ModuleRunner, "消息模板组变量", tv.Key, "无法识别!")
					return SUCCESSFUL
				}
			}
			for b := range x.Templates {
				vs = append(vs, &x.Templates[b])
			}
		} else if tv.IsFlow {
			c := &t.Config
			Info(ModuleRunner, "识别流程变量", tv.ID)
			ft = GetFlow(tv.Node, c, st)
			if ft == nil {
				Err(ModuleRunner, "无法识别流程ID", tv.Node, "!")
				return SUCCESSFUL
			}
			Info(ModuleRunner, "加载流程变量", tv.ID, "流程ID", ft.FID)
			fd = fds[ft.FID]
			// 复制用例的上下文
			for k, ctx := range cd.Ctx {
				h := strings.Split(k, DOT)
				if len(h) < 1 {
					continue
				}
				// 只复制全局(x)的上下文
				if h[0] != CTX {
					continue
				}
				fd.Ctx[k] = ctx
				if k == CTX {
					_, ok := fd.Msgs[CTX]
					if !ok {
						fd.Msgs[CTX] = ctx
					} else {
						fd.Msgs[CTX] = UpdateJSON(fd.Msgs[CTX], ctx)
					}
				} else {
					_, ok := fd.Vars[k]
					if !ok {
						fd.Vars[k] = ctx
					} else {
						fd.Vars[k] = UpdateJSON(fd.Vars[k], ctx)
					}
				}
			}
			// 复制用例的流程上下文
			for k, ctx := range cd.Ctx {
				h := strings.Split(k, DOT)
				if len(h) < 1 {
					continue
				}
				// 只复制本流程变量上下文
				if h[0] != tv.Key {
					continue
				}
				h[0] = CTX
				v := strings.Join(h, DOT)
				fd.Ctx[v] = ctx
				_, ok := fd.Vars[v]
				if !ok {
					fd.Vars[v] = ctx
				} else {
					fd.Vars[v] = UpdateJSON(fd.Vars[k], ctx)
				}
			}
			tp := getTemplates(ft.Vars)
			for total > 0 {
				total--
				runFlowTemplates(id, key, tv.Key, r, t, st, ct, pc, ft, op, cd, pd, fd, fds, tp)
			}
			// 更新用例的上下文
			if ct.Local == 1 {
				continue
			}
			for k, ctx := range fd.Ctx {
				h := strings.Split(k, DOT)
				// 只复制全局(x)的上下文
				if len(h) < 1 || h[0] != CTX {
					continue
				}
				cd.Ctx[k] = ctx
			}
			continue
		} else {
			vs = append(vs, tv)
		}
		// 循环次数，不是消息模板组时，total为1
		for i := 0; i < total; i++ {
			for _, k := range vs {
				for j := 0; j < k.Count; j++ {
					if k.IsWait {
						k, _ := strconv.Atoi(k.Key)
						time.Sleep(time.Duration(k) * time.Millisecond)
						break
					}
					switch key {
					case REQUEST:
						if SUCCESSFUL != Request(id, r, t, ct, pc, nil, op, k, cd, pd, nil) {
							return FAILED
						}
					}
				}
			}
		}
	}
	return SUCCESSFUL
}

func runFlowTemplates(id int, key, vk string, r TaskRequest, t *Testing, st *TestSuite,
	ct *CaseType, pc *ControlType, ft *FlowType, op *OpType, cd, pd, fd *VarData,
	fds map[string]*VarData, p *VarType) int {
	var vs []*TemplateVarType
	n := len(p.Templates)
	if ft.Random > n {
		n = ft.Random
	}
	for n > 0 {
		// 取指针
		vs = nil
		a := rand.Intn(len(p.Templates))
		if ft.Random == 0 {
			a = len(p.Templates) - n
		}
		tv := &p.Templates[a]
		n--
		total := 1
		// 消息模板组或流程
		if tv.IsVar {
			total = tv.Count
			x, ok := ft.Vars[tv.Key]
			if !ok {
				Err(ModuleRunner, "消息模板组变量", tv.Key, "无法识别!")
				return SUCCESSFUL
			}
			for b := range x.Templates {
				vs = append(vs, &x.Templates[b])
			}
		} else if tv.IsFlow {
			c := &t.Config
			Info(ModuleRunner, "识别流程变量", tv.ID)
			fp := GetFlow(tv.Node, c, st)
			if fp == nil {
				Err(ModuleRunner, "无法识别流程ID", tv.Node, "!")
				return SUCCESSFUL
			}
			Info(ModuleRunner, "流程变量", tv.ID, "流程ID", fp.FID)
			fdd := fds[fp.FID]
			// 复制用例的全局上下文
			for k, ctx := range cd.Ctx {
				h := strings.Split(k, DOT)
				if len(h) < 1 {
					continue
				}
				// 只复制全局(x)的上下文
				if h[0] != CTX {
					continue
				}
				fdd.Ctx[k] = ctx
				if k == CTX {
					_, ok := fdd.Msgs[CTX]
					if !ok {
						fdd.Msgs[CTX] = ctx
					} else {
						fdd.Msgs[CTX] = UpdateJSON(fdd.Msgs[CTX], ctx)
					}
				} else {
					_, ok := fdd.Vars[k]
					if !ok {
						fdd.Vars[k] = ctx
					} else {
						fdd.Vars[k] = UpdateJSON(fdd.Vars[k], ctx)
					}
				}
			}
			// 复制用例的流程上下文
			for k, ctx := range cd.Ctx {
				h := strings.Split(k, DOT)
				if len(h) < 1 {
					continue
				}
				// 只复制本流程变量上下文
				if h[0] != vk {
					continue
				}
				h[0] = CTX
				v := strings.Join(h, DOT)
				fdd.Ctx[v] = ctx
				_, ok := fdd.Vars[v]
				if !ok {
					fdd.Vars[v] = ctx
				} else {
					fdd.Vars[v] = UpdateJSON(fdd.Vars[k], ctx)
				}
			}
			tp := getTemplates(fp.Vars)
			for total > 0 {
				total--
				runFlowTemplates(id, key, tv.Key, r, t, st, ct, pc, fp, op, cd, pd, fdd, fds, tp)
			}
			if fp.Local == 1 {
				continue
			}
			// 更新用例的上下文
			for k, ctx := range fdd.Ctx {
				h := strings.Split(k, DOT)
				// 只复制全局(x)的上下文
				if len(h) < 1 || h[0] != CTX {
					continue
				}
				cd.Ctx[k] = ctx
			}
			continue
		} else {
			vs = append(vs, tv)
		}
		// 循环次数，不是消息模板组时，total为1
		for i := 0; i < total; i++ {
			for _, k := range vs {
				for j := 0; j < k.Count; j++ {
					if k.IsWait {
						k, _ := strconv.Atoi(k.Key)
						time.Sleep(time.Duration(k) * time.Millisecond)
						break
					}
					switch key {
					case REQUEST:
						if SUCCESSFUL != Request(id, r, t, ct, pc, ft, op, k, cd, pd, fd) {
							return FAILED
						}
					}
				}
			}
		}
	}
	return SUCCESSFUL
}

func isArg(key string, args map[string]int) bool {
	_, ok := args[key]
	return ok
}

func getTemplates(p map[string]*VarType) *VarType {
	for _, v := range p {
		Info(ModuleRunner, "获取模板组：变量", v.Key, "是否top", v.IsTop, "是否模板", v.IsTemplate, "模板数", len(v.Templates))
		if v.IsTop && v.IsTemplate && len(v.Templates) > 0 && !v.IsGroup {
			// 找到第一个节点
			for _, k := range v.Templates {
				Info(ModuleRunner, "遍历模板组：变量", k.Key, "ID", k.ID, "是否节点", k.IsNode,
					"是否变量", k.IsVar, "是否等待", k.IsWait, "是否上下文", k.IsCtx, "目标列表", k.Targets, "重复次数", k.Count, "模板数", len(v.Templates))
				if k.IsNode {
					return v
				}
			}
		}
	}
	return nil
}
