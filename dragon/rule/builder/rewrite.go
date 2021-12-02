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

package builder

import (
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func isDefined(key string, m map[string][]string) bool {
	_, ok := m[key]
	return ok
}

func refineRangeValue(s []string) []string {
	var v []string
	for _, k := range s {
		t := strings.Split(k, OpDotDot)
		if strings.Index(k, OpDotDotDot) > 0 || len(t) < 2 || strings.Index(k, COLON) > 0 {
			v = append(v, k)
			continue
		}
		start, e1 := strconv.Atoi(t[0])
		end, e2 := strconv.Atoi(t[1])
		if e1 != nil || e2 != nil {
			continue
		}
		if end-start > 2*MaxRange {
			for i := start; i < start+MaxRange; i++ {
				v = append(v, strconv.Itoa(i))
			}
			for i := end; i > end-MaxRange; i-- {
				v = append(v, strconv.Itoa(i))
			}
			x := (end - start) / 2
			v = append(v, strconv.Itoa(x))
			continue
		}
		for i := start; i <= end; i++ {
			v = append(v, strconv.Itoa(i))
		}
	}
	return v
}

func convertOp(r TaskRequest, meta map[string]int, rs map[string][]string, ops []string) []string {
	var v []string
	for _, value := range refineRangeValue(ops) {
		if isDefined(value, rs) {
			for _, k := range refineRangeValue(rs[value]) {
				meta[k] = 1
				v = append(v, k)
			}
		} else {
			v = append(v, value)
		}
	}
	Info(ModuleBuilder, "变换", ops, "为", v)
	return v
}

func buildArrayData(p *ConstructType, t *OpType) []OpType {
	n := 1
	for i := range p.Settings {
		op := p.Settings[i]
		if len(op.Value) > n {
			n = len(op.Value)
		}
	}
	ds := make([]OpType, n)
	for i := range ds {
		op := &ds[i]
		op.Key = t.Key
		op.Op = t.Op
		for _, v := range p.Formats {
			op.Value = append(op.Value, v.Content)
		}
		if len(p.Formats) == 0 {
			op.Value = []string{EMPTY}
		}
	}
	for i := range p.Settings {
		op := p.Settings[i]
		m := len(op.Value)
		r := rand.Intn(m) // 起始位置随机
		for k := range ds {
			v := op.Value[(k+r)%m]
			s := &ds[k]
			for j := range s.Value {
				if len(p.Formats) == 0 {
					h := []string{s.Value[j], v}
					s.Value[j] = strings.Join(h, EMPTY)
					continue
				}
				s.Value[j] = strings.Replace(s.Value[j], op.Key, v, UNDEFINED)
			}
		}
	}
	if p.IsInt {
		for k := range ds {
			s := &ds[k]
			for j := range s.Value {
				s.Value[j] = Calculate(s.Value[j])
			}
		}
	}
	return ds
}

func buildStringData(p *ConstructType, t *OpType) []OpType {
	n := 1
	for i := range p.Settings {
		op := p.Settings[i]
		if len(op.Value) > n {
			n = len(op.Value)
		}
	}
	ds := make([]OpType, n)
	for i := range ds {
		op := &ds[i]
		op.Key = t.Key
		op.Op = t.Op
		op.Value = append(op.Value, p.Formats[0].Content)
	}
	for i := range p.Settings {
		op := p.Settings[i]
		m := len(op.Value)
		r := rand.Intn(m) // 起始位置随机
		for k := range ds {
			s := &ds[k]
			v := op.Value[(k+r)%m]
			s.Value[0] = strings.Replace(s.Value[0], op.Key, v, UNDEFINED)
		}
	}
	if p.IsInt {
		for k := range ds {
			s := &ds[k]
			s.Value[0] = Calculate(s.Value[0])
		}
	}
	return ds
}

func buildData(p *ConstructType, t *OpType) []OpType {
	if p.IsList {
		return buildArrayData(p, t)
	}
	return buildStringData(p, t)
}

func fetchData(r TaskRequest, d map[string]*ConstructType, p *SettingType, m map[string]*VarType) {
	var ds []OpType
	found := 0
	p.Items = 0
	for i := range p.Ops {
		t := &p.Ops[i]
		for _, v := range t.Value {
			k, ok := d[v]
			if ok {
				ot := buildData(k, t)
				ds = append(ds, ot...)
				Info(ModuleBuilder, "映射数据构造", v, "为", ot)
				found = 1
				if len(k.Formats) > 0 {
					p.MustString = k.Formats[0].MustString
				}
				if k.IsList {
					p.IsList = k.IsList
					n, e := strconv.Atoi(k.Keys[0].Content)
					if e == nil {
						p.Items += n
					}
				}
				if k.IsInt {
					p.IsInt = k.IsInt
				}
			}
			a := strings.Split(t.Key, DOT)
			b := strings.Split(v, DOT)
			_, ok1 := m[a[0]]
			_, ok2 := m[b[0]]
			if ok1 && ok2 {
				p.IsDep = true
				Info(ModuleBuilder, "命中依赖赋值", t.Key, v)
			} else if ok1 {
				p.IsField = true
				Info(ModuleBuilder, "命中字段赋值", t.Key, v)
			} else {
				// 动态标签须以__开始，以__结尾，如 __tagA__
				if StrMatchStart(a[0], DoubleUnderScope) && StrMatchEnd(a[0], DoubleUnderScope) {
					p.IsTag = true
					Info(ModuleBuilder, "命中标签赋值", t.Key, v)
				}
			}
		}
	}
	if found == 1 {
		Info(ModuleBuilder, "更新赋值定义前", p.Ops)
		p.Ops = []OpType{}
		p.Ops = append(p.Ops, ds...)
		Info(ModuleBuilder, "更新赋值定义后", p.Ops)
	}
}

func fetchMessageData(r TaskRequest, d map[string]*ConstructType, p *ConstructType) {
	var ds []OpType
	found := 0
	for i := range p.Settings {
		t := &p.Settings[i]
		b := OpType{Key: t.Key, Op: t.Op}
		for _, v := range t.Value {
			k, ok := d[v]
			if ok {
				ot := buildData(k, t)
				for _, s := range ot {
					b.Value = append(b.Value, s.Value...)
				}
				Debug(ModuleBuilder, "映射数据构造", v, "为", ot)
				found = 1
			} else {
				b.Value = append(b.Value, v)
			}
		}
		ds = append(ds, b)
	}
	if found == 1 {
		Debug(ModuleBuilder, "更新赋值定义前", p.Settings)
		p.Settings = []OpType{}
		p.Settings = append(p.Settings, ds...)
		Debug(ModuleBuilder, "更新赋值定义后", p.Settings)
	}
}

func fetchVars(r TaskRequest, key string, i int, v map[string]*VarType,
	p *VarType, op *OpType, sp *TestSuite, c *ConfigType) int {
	Debug(ModuleBuilder, "变量前缀", key, op.Value)
	id := i
	for _, value := range op.Value {
		s := strconv.Itoa(id + 1)
		h := []string{key, s}
		a := strings.Join(h, EMPTY)
		Debug(ModuleBuilder, "识别变量", value)
		t := strings.Split(value, OpColon)
		n := len(t)
		if n > 1 {
			k := t[0]
			b := strings.Join(t[1:n], OpColon)
			vt := TemplateVarType{ID: k, Key: b, Count: 1}
			_, ok := v[b]
			if ok {
				vt.IsVar = true
			}
			if c != nil {
				x := t[1]
				if strings.Contains(x, STAR) {
					h = strings.Split(x, STAR)
					x = h[1]
				}
				h = strings.Split(x, DOT)
				sk := SearchLongestMatchNode(h, c.Nodes, DOT)
				if sk != EMPTY {
					vt.Node = sk
					vt.IsNode = true
				}
				if !vt.IsNode {
					_, err := strconv.Atoi(vt.Key)
					switch {
					case err == nil:
						vt.IsWait = true
					case IsFlow(x, c, sp):
						vt.IsFlow = true
						vt.Node = x
					case IsGroup(x, v):
						v[x].IsGroup = true
						vt.IsGroup = true
						vt.IsVar = true
					default:
						vt.IsCtx = true
					}
				}
			}
			if vt.IsNode || vt.IsFlow {
				v[k] = &VarType{Key: k, Node: vt.Node, IsTemplate: true, IsNode: vt.IsNode, IsFlow: vt.IsFlow, Templates: []TemplateVarType{vt}}
				v[a] = &VarType{Key: a, Node: vt.Node, IsTemplate: true, IsNode: vt.IsNode, IsFlow: vt.IsFlow, Templates: []TemplateVarType{vt}}
			} else if vt.IsCtx {
				_, ok := v[k]
				if !ok {
					v[k] = &VarType{Key: k, IsTemplate: true, Templates: []TemplateVarType{vt}}
				} else {
					v[k].Templates = append(v[k].Templates, vt)
				}
			}
			p.Templates = append(p.Templates, vt)
			Debug(ModuleBuilder, "生成变量", k, "值", b, "变量前缀", key, "是否变量", vt.IsVar, "是否节点", vt.IsNode, "是否流程", vt.IsFlow, "是否等待", vt.IsWait, "是否上下文", vt.IsCtx, "模板列表", p.Templates)
			Debug(ModuleBuilder, "生成变量", a, "值", b, "变量前缀", key, "是否变量", vt.IsVar, "是否节点", vt.IsNode, "是否流程", vt.IsFlow, "是否等待", vt.IsWait, "是否上下文", vt.IsCtx, "模板列表", p.Templates)
		} else {
			vt := TemplateVarType{ID: a, Key: value, Count: 1}
			_, ok := v[value]
			if ok {
				vt.IsVar = true
			}
			if c != nil {
				if strings.Contains(value, STAR) {
					h = strings.Split(value, STAR)
					value = h[1]
				}
				h = strings.Split(value, DOT)
				sk := SearchLongestMatchNode(h, c.Nodes, DOT)
				if sk != EMPTY {
					vt.Node = sk
					vt.IsNode = true
				}
				if !vt.IsNode {
					_, err := strconv.Atoi(vt.Key)
					switch {
					case err == nil:
						vt.IsWait = true
					case IsFlow(value, c, sp):
						vt.IsFlow = true
						vt.Node = value
					case IsGroup(value, v):
						v[value].IsGroup = true
						vt.IsGroup = true
						vt.IsVar = true
					default:
						vt.IsCtx = true
						vt.ID = CTX
					}
				}
			}
			if vt.IsNode || vt.IsFlow {
				v[a] = &VarType{Key: a, Node: vt.Node, IsTemplate: true, IsNode: vt.IsNode, IsFlow: vt.IsFlow, Templates: []TemplateVarType{vt}}
			} else if vt.IsCtx {
				_, ok := v[CTX]
				if !ok {
					v[CTX] = &VarType{Key: CTX, IsTemplate: true, Templates: []TemplateVarType{vt}}
				} else {
					v[CTX].Templates = append(v[CTX].Templates, vt)
				}
			}
			p.Templates = append(p.Templates, vt)
			Debug(ModuleBuilder, "生成变量", a, "值", value, "变量前缀", key, "是否变量", vt.IsVar, "是否上下文", vt.IsCtx,
				"是否节点", vt.IsNode, "是否流程", vt.IsFlow, "是否等待", vt.IsWait, "模板列表", p.Templates)
		}
		id += 1
	}
	return id
}

func isIntType(v []string) bool {
	var e error
	for i := range v {
		_, e = strconv.Atoi(v[i])
		if e != nil {
			return false
		}
	}
	return true
}

func convertConstructor(r TaskRequest, meta map[string]int, rs map[string][]string, p *ConstructType) {
	isInt := true
	for i := range p.Settings {
		t := &p.Settings[i]
		t.Value = convertOp(r, meta, rs, t.Value)
		if !isIntType(t.Value) {
			isInt = false
		}
	}
	p.IsInt = isInt
}

func convertTestCase(r TaskRequest, rs map[string][]string, m, d map[string]*ConstructType,
	t map[string]*TemplateType, p *CaseType, sp *TestSuite, c *ConfigType) {
	for i := range p.Steps {
		step := &p.Steps[i]
		convertStep(r, rs, m, d, t, p.Vars, step, sp, c)
		for j := range step.Expects {
			op := &step.Expects[j]
			op.Value = convertOp(r, c.Meta, rs, op.Value)
			Debug(ModuleBuilder, "期望内容", op.Key, op.Value)
		}
	}
	mf := map[string]int{}
	for _, v := range p.Vars {
		Info(ModuleBuilder, "变量", v.Key, "是否模板", v.IsTemplate, "是否操作", v.IsOpt,
			"是否顶级变量", v.IsTop, "是否节点", v.IsNode, "是否流程", v.IsFlow, "模板列表", v.Templates)
		if v.IsFlow {
			_, ok := mf[v.Key]
			if !ok {
				p.FIDs = append(p.FIDs, v.Key)
			}
		}
	}
	for k, v := range p.Vars {
		if v.IsTemplate {
			rewriteTemplates(r, k, t, v, p.Vars)
			Info(ModuleBuilder, "构造变量", k, p.Vars[k])
		}
	}
}

func convertStep(r TaskRequest, rs map[string][]string, m, d map[string]*ConstructType,
	t map[string]*TemplateType, p map[string]*VarType, step *StepType, sp *TestSuite, c *ConfigType) {
	var vp *VarType
	vp = nil
	id := 0
	for j := range step.Setting.Ops {
		op := &step.Setting.Ops[j]
		op.Value = convertOp(r, c.Meta, rs, op.Value)
		if vp != nil && op.Op != EMPTY {
			vp.IsTemplate = false
			vp.IsOpt = true
			tp := &step.Setting.Ops[0]
			// 非模板，非字典变量
			if vp.Value == nil && len(tp.Value) > 0 {
				setVar(vp, tp)
			}
			// 字典变量
			if op.Op == OpEqual && len(op.Value) > 0 {
				s := []string{vp.Key, op.Key}
				key := strings.Join(s, DOT)
				p[key] = &VarType{Key: key, IsOpt: true, Op: op.Op}
				p[key].Value = []interface{}{}
				setVar(p[key], op)
				continue
			} else if op.Op != EMPTY {
				break
			}
		}
		if op.Op == OpConstruct || op.Op == OpEqual {
			if strings.Index(op.Key, OpDot) == UNDEFINED && strings.Index(op.Key, DoubleUnderScope) == UNDEFINED {
				p[op.Key] = &VarType{Key: op.Key}
				vp = p[op.Key]
				if op.Op == OpConstruct {
					vp.IsTemplate = true
				} else {
					vp.IsOpt = true
				}
				if id == 0 {
					vp.IsTop = true
					vp.Op = op.Op // 只有->操作符才是构造变量，遍历取值
				}
			}
		}
		if vp != nil && vp.IsTemplate && len(op.Value) > 0 {
			id = fetchVars(r, vp.Key, id, p, vp, op, sp, c)
		}
	}
	fetchData(r, d, &step.Setting, p)
}

func setVar(vp *VarType, op *OpType) {
	if op.IsMap {
		vp.Value = append(vp.Value, map[string]interface{}{})
	}
	if op.IsList {
		vp.Value = append(vp.Value, op.Value)
	}
	if isIntType(op.Value) {
		for _, v := range op.Value {
			k, _ := strconv.Atoi(v)
			vp.Value = append(vp.Value, k)
		}
	} else {
		for _, v := range op.Value {
			vp.Value = append(vp.Value, v)
		}
	}
}

func rewriteTemplates(r TaskRequest, key string, t map[string]*TemplateType, v *VarType, p map[string]*VarType) {
	var ds []TemplateVarType
	expand := false
	for _, k := range v.Templates {
		Debug(ModuleBuilder, "加载模板前", key, k)
		h := strings.Split(k.Key, STAR)
		n := len(h)
		if n > 1 {
			c, e := strconv.Atoi(h[0])
			if e == nil {
				k.Count = c
				k.Key = strings.Join(h[1:n], WILDCARD)
			} else {
				k.Key = strings.Join(h, WILDCARD)
			}
			expand = true
		}
		h = strings.Split(k.Key, COLON)
		n = len(h)
		if n > 1 {
			k.Key = h[0]
			a := strings.Split(h[1], DOT)
			n = len(a)
			b := false
			for _, s := range a {
				i, e := strconv.Atoi(s)
				if e == nil {
					if b {
						m := len(k.Targets)
						if m > 0 && k.Targets[m-1] < i {
							for x := k.Targets[m-1] + 1; x <= i; x++ {
								k.Targets = append(k.Targets, x)
							}
						}
						b = false
						continue
					}
					k.Targets = append(k.Targets, i)
				}
				if s == EMPTY {
					b = true
				}
			}
			expand = true
		}
		Debug(ModuleBuilder, "加载模板后", key, k)
		if strings.Index(k.Key, WILDCARD) == UNDEFINED {
			ds = append(ds, k)
			continue
		}
		reg := regexp.MustCompile(k.Key)
		j := strings.Index(k.Key, string(os.PathSeparator))
		prefix := GetPath(TPath, DefaultTPath)
		for s := range t {
			i := strings.Index(s, string(os.PathSeparator))
			if j == UNDEFINED && i != UNDEFINED {
				continue
			} else if j != UNDEFINED && i == UNDEFINED {
				continue
			}
			if reg.MatchString(s) {
				if strings.Index(s, prefix) == 0 {
					continue
				}
				d := TemplateVarType{ID: k.ID, Key: s, Count: k.Count, IsNode: k.IsNode}
				for _, i := range k.Targets {
					d.Targets = append(d.Targets, i)
				}
				ds = append(ds, d)
				expand = true
			}
		}
	}
	if expand {
		p[key] = &VarType{Key: v.Key, IsTemplate: v.IsTemplate, IsFlow: v.IsFlow,
			IsTop: v.IsTop, IsOpt: v.IsOpt, IsNode: v.IsNode, Templates: ds}
		Log(ModuleBuilder, "扩展模板后", key, p[key].Templates)
	}
	Log(ModuleBuilder, "最终模板", key, p[key].Templates)
}

func convertControlFlow(r TaskRequest, rs map[string][]string, m, d map[string]*ConstructType,
	t map[string]*TemplateType, p *ControlType, c *ConfigType) {
	for i := range p.Flows {
		step := &p.Flows[i]
		convertStep(r, rs, m, d, t, p.Vars, step, nil, c)
	}
	for _, v := range p.Vars {
		Info(ModuleBuilder, "变量", v.Key, "是否模板", v.IsTemplate, "是否操作", v.IsOpt,
			"是否顶级变量", v.IsTop, "模板列表", v.Templates)
	}
	for k, v := range p.Vars {
		if v.IsTemplate {
			rewriteTemplates(r, k, t, v, p.Vars)
			Info(ModuleBuilder, "构造变量", k, p.Vars[k])
		}
	}
}

func convertFlow(r TaskRequest, rs map[string][]string, m, d map[string]*ConstructType,
	t map[string]*TemplateType, p *FlowType, sp *TestSuite, c *ConfigType) {
	for i := range p.Steps {
		step := &p.Steps[i]
		convertStep(r, rs, m, d, t, p.Vars, step, sp, c)
	}
	mf := map[string]int{}
	for _, v := range p.Vars {
		Info(ModuleBuilder, "变量", v.Key, "是否模板", v.IsTemplate, "是否操作", v.IsOpt,
			"是否顶级变量", v.IsTop, "是否节点", v.IsNode, "是否流程", v.IsFlow, "模板列表", v.Templates)
		if v.IsFlow {
			_, ok := mf[v.Key]
			if !ok {
				p.FIDs = append(p.FIDs, v.Key)
			}
		}
	}
	for k, v := range p.Vars {
		if v.IsTemplate {
			rewriteTemplates(r, k, t, v, p.Vars)
			Info(ModuleBuilder, "构造变量", k, p.Vars[k])
		}
	}
}

func convertNode(r TaskRequest, meta map[string]int, rs map[string][]string, m, d map[string]*ConstructType, p *NodeType) {
	k := map[string]*VarType{}
	for i := range p.Steps {
		step := &p.Steps[i]
		for j := range step.Setting.Ops {
			op := &step.Setting.Ops[j]
			op.Value = convertOp(r, meta, rs, op.Value)
		}
		fetchData(r, d, &step.Setting, k)
	}
}

func convertSuite(r TaskRequest, rs map[string][]string, m, d map[string]*ConstructType,
	t map[string]*TemplateType, p *TestSuite, c *ConfigType) {
	for k := range p.ControlFlow {
		convertControlFlow(r, rs, m, d, t, p.ControlFlow[k], c)
	}
	for k := range p.Setup {
		convertControlFlow(r, rs, m, d, t, p.Setup[k], c)
	}
	for k := range p.Teardown {
		convertControlFlow(r, rs, m, d, t, p.Teardown[k], c)
	}
	for k := range p.Checks {
		convertControlFlow(r, rs, m, d, t, p.Checks[k], c)
	}
	for k := range p.Flows {
		convertFlow(r, rs, m, d, t, p.Flows[k], p, c)
	}
	for _, g := range p.Groups {
		for _, v := range g.Cases {
			convertTestCase(r, rs, m, d, t, v, p, c)
		}
	}
}

func rewrite(r TaskRequest, p *Testing) {
	rand.Seed(time.Now().Unix())
	c := &p.Config
	rs := c.Settings
	m := c.Messages
	d := c.Data
	for k := range c.Data {
		convertConstructor(r, c.Meta, rs, c.Data[k])
	}
	for k := range c.Messages {
		convertConstructor(r, c.Meta, rs, c.Messages[k])
		fetchMessageData(r, c.Data, c.Messages[k])
	}
	for k := range c.Nodes {
		convertNode(r, c.Meta, rs, m, d, c.Nodes[k])
	}
	for k := range c.Flows {
		convertFlow(r, rs, m, d, p.Templates, c.Flows[k], nil, c)
	}
	for k := range p.ControlFlow {
		convertControlFlow(r, rs, m, d, p.Templates, p.ControlFlow[k], c)
	}
	for k := range p.Setup {
		convertControlFlow(r, rs, m, d, p.Templates, p.Setup[k], c)
	}
	for k := range p.Teardown {
		convertControlFlow(r, rs, m, d, p.Templates, p.Teardown[k], c)
	}
	for k := range p.Checks {
		convertControlFlow(r, rs, m, d, p.Templates, p.Checks[k], c)
	}
	for k := range p.Suites {
		convertSuite(r, rs, m, d, p.Templates, p.Suites[k], c)
	}
}
