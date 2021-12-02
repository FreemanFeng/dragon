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

package rule

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func UpdateVar(id, vt int, proto string, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType, ft *FlowType,
	tp *TemplateVarType, vd *VarData) int {
	// 先映射消息变量
	for _, v := range vd.Data[tp.ID] {
		s := []string{tp.ID, v.Top}
		key := strings.Join(s, DOT)
		if v.Top != MainMsg {
			vd.Tops[tp.ID] = append(vd.Tops[tp.ID], key)
		}
	}
	// 需要加载变量对应的消息模板
	for _, v := range vd.Data[tp.ID] {
		Log(ModuleRunner, "模板路径", v.Path, "名称", v.Name, "顶级目录", v.Top, "变量", tp.ID, "是否JSON", v.IsJSON)
		if SUCCESSFUL != updateTags(vt, r, t, ct, pc, vd, v) {
			Err(ModuleRule, "!!!!!!更新标签变量失败!")
			return FAILED
		}
		if v.IsJSON {
			msg := ToJSON(v.Content)
			if msg == nil {
				Err(ModuleRule, "!!!!!!无法生成JSON数据!")
				return FAILED
			}
			s := []string{tp.ID, v.Top}
			key := strings.Join(s, DOT)
			sn := strings.Split(v.Name, DOT)
			sk := SearchLongestMatchNode(sn, t.Config.Nodes, DOT)
			service := t.Config.Nodes[sk]
			if v.Top == MainMsg || tp.ID == CTX || sk == EMPTY {
				key = tp.ID
			}
			_, ok := vd.Msgs[key]
			if key != CTX || !ok {
				vd.Msgs[key] = msg
			} else {
				// 全局上下文，可多次叠加
				vd.Msgs[key] = UpdateJSON(vd.Msgs[key], msg)
			}
			switch vt {
			case CaseVar:
				if sk != EMPTY {
					for i := range service.Steps {
						sp := &service.Steps[i]
						updateValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					}
				}
				for i := range ct.Steps {
					sp := &ct.Steps[i]
					updateValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
				}
				if sk != EMPTY {
					for i := range service.Steps {
						sp := &service.Steps[i]
						updateDepValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
						// 带条件on的赋值，要满足条件
						updateMessage(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					}
				}
				for i := range ct.Steps {
					sp := &ct.Steps[i]
					updateDepValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					// 带条件on的赋值，要满足条件
					updateMessage(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
				}
			case FlowVar:
				if sk != EMPTY {
					for i := range service.Steps {
						sp := &service.Steps[i]
						updateValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					}
				}
				for i := range ft.Steps {
					sp := &ft.Steps[i]
					updateValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
				}
				if sk != EMPTY {
					for i := range service.Steps {
						sp := &service.Steps[i]
						updateDepValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
						// 带条件on的赋值，要满足条件
						updateMessage(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					}
				}
				for i := range ft.Steps {
					sp := &ft.Steps[i]
					updateDepValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					// 带条件on的赋值，要满足条件
					updateMessage(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
				}
			case CFVar:
				if sk != EMPTY {
					for i := range service.Steps {
						sp := &service.Steps[i]
						updateValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
						updateDepValue(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
						// 带条件on的赋值，要满足条件
						updateMessage(id, vt, tp.ID, key, msg, r, t, ct, pc, v, vd, ct.Vars, sp)
					}
				}
				for i := range pc.Flows {
					sp := &pc.Flows[i]
					updateValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, pc.Vars, sp)
					updateDepValue(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, pc.Vars, sp)
					// 带条件on的赋值，要满足条件
					updateMessage(id, vt, EMPTY, key, msg, r, t, ct, pc, v, vd, pc.Vars, sp)
				}
			}
		}
	}
	// 将消息序列化后回写消息原始内容
	for _, v := range vd.Data[tp.ID] {
		Log(ModuleRunner, "回写消息模板路径", v.Path, "名称", v.Name, "顶级目录", v.Top, "是否JSON", v.IsJSON)
		if v.IsJSON {
			s := []string{tp.ID, v.Top}
			key := strings.Join(s, DOT)
			if v.Top == MainMsg {
				key = tp.ID
			}
			msg, ok := vd.Msgs[key]
			if ok {
				ts := []string{key, RAFile}
				fkey := strings.Join(ts, DOT)
				ts = []string{key, RAFileZ}
				fkey2 := strings.Join(ts, DOT)
				b, e := json.Marshal(msg)
				if e == nil {
					v.Content = b
					vd.Vars[fkey] = b
					vd.Vars[fkey2] = b
				}
			}
		}
	}
	return SUCCESSFUL
}

func updateTags(vt int, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType, p *VarData, v *FileType) int {
	b := v.Content
	for _, tag := range v.Tags {
		op, ok := p.Tags[tag]
		if !ok {
			Err(ModuleRule, "文件", v.Name, "找不到标签", tag)
			continue
		}
		// 函数调用
		if len(op.Value) > 0 && !IsConfig(op.Value[0], t) && !IsDigit(op.Value[0]) && IsCall(op.Value[0], t) {
			key := op.Value[0]
			RunBuiltin(vt, EMPTY, EMPTY, key, r, t, ct, pc, op, p)
			RunFuncs(vt, EMPTY, EMPTY, key, r, t, ct, pc, op, p)
			if SUCCESSFUL != RunBins(vt, EMPTY, EMPTY, key, r, t, ct, pc, op, p) {
				return FAILED
			}
			k, ok := p.Vars[key]
			if !ok {
				Err(ModuleRule, "无法调用", key, "赋值给", op.Key)
				return FAILED
			}
			v.Content = bytes.Replace(b, []byte(tag), []byte(IfToString(k)), UNDEFINED)
			continue
		}
		n := len(op.Value)
		_, ok = p.VPos[tag]
		if !ok {
			p.VPos[tag] = UNDEFINED
		}
		switch ct.Fill {
		case FillAll:
			p.VPos[tag] = (p.VPos[tag] + 1) % n
		case FillRandom:
			p.VPos[tag] = rand.Intn(n)
		}
		i := p.VPos[tag]
		k := op.Value[i]
		x, ok := p.Vars[k]
		if !ok {
			// 非变量，普通替换
			b = bytes.Replace(b, []byte(tag), []byte(k), UNDEFINED)
		} else {
			// 变量，用变量值赋值
			b = bytes.Replace(b, []byte(tag), []byte(IfToString(x)), UNDEFINED)
		}
	}
	v.Content = b
	return SUCCESSFUL
}

func updateComposingKeys(id, vt, n int, key, v string, msg interface{}, vd *VarData) {
	if !IsOpVars(v) && StrMatchStart(v, AddDot(key)...) {
		k, ok := vd.Keys[v]
		if !ok {
			vd.Keys[v] = BuildKeys(v)
			k = vd.Keys[v]
		}
		updateWildcardKeys(id, vt, n, v, msg, k, vd)
	} else if IsOpVars(v) {
		for _, vv := range SplitVars(v) {
			if IsOpVars(vv) || !StrMatchStart(vv, AddDot(key)...) {
				continue
			}
			k, ok := vd.Keys[vv]
			if !ok {
				vd.Keys[vv] = BuildKeys(vv)
				k = vd.Keys[vv]
			}
			updateWildcardKeys(id, vt, n, vv, msg, k, vd)
		}
	}
}
func updateWildcardKeys(id, vt, n int, key string, msg interface{}, k []FormatType, vd *VarData) {
	Info(ModuleRule, ">>>>> 更新Key", key)
	v, ok := vd.Vars[key]
	if ok {
		Info(ModuleRule, ">>>>> Key", key, "已存在值", v)
		return
	}
	if strings.Contains(key, OpDotDotDot) {
		m := map[string]interface{}{}
		GetAllJSON(msg, k[n:], m)
		if len(m) > 0 {
			vd.Vars[key] = m
			Info(ModuleRule, ">>>>> 更新模糊key", key, "更新值为", m, "类型为", TypeOf(m))
		}
	} else if strings.Contains(key, DOT) {
		x := GetJSON(msg, k[n:])
		if x != nil {
			vd.Vars[key] = x
			Info(ModuleRule, ">>>>> 更新key", key, "更新值为", x, "类型为", TypeOf(x))
		}
	}
}

// 更新字段值
func updateValue(id, vt int, mid, key string, msg interface{}, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType,
	fp *FileType, vd *VarData, vp map[string]*VarType, sp *StepType) int {
	var vr []interface{}
	s := strings.Split(key, DOT)
	ts := []string{key, RAFile}
	fkey := strings.Join(ts, DOT)
	ts = []string{key, RAFileZ}
	fkey2 := strings.Join(ts, DOT)
	n := len(s)
	kon := SearchOp(OpOn, sp.Setting.Ops)
	// 更新赋值区域字段变量
	for i := range sp.Setting.Ops {
		// 跳过条件区
		p := &sp.Setting.Ops[i]
		if !FoundID(id, p.CIDs) {
			return SUCCESSFUL
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]

		if ok && tv.IsTemplate {
			continue
		}
		for _, v := range p.Value {
			if v == strings.TrimSpace(NoChars) {
				continue
			}
			if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
				v = strings.Join([]string{mid, v}, DOT)
				LoadKeys(vd, vp, v)
			}
			// 其他消息变量字段
			if fp.Top == MainMsg && StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				continue
			}
			if Equal(v, fkey, fkey2) {
				Log(ModuleRule, "获取文件", fp.Name, "内容存到变量", v)
				vd.Vars[v] = fp.Content
				vd.Vars[fkey] = fp.Content
				vd.Vars[fkey2] = fp.Content
				continue
			}
			updateComposingKeys(id, vt, n, key, v, msg, vd)
		}
		// 其他消息变量字段
		if fp.Top == MainMsg && StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
			continue
		}
		// 更新op.Key
		if StrMatchStart(sk, AddDot(key)...) {
			k, ok := vd.Keys[sk]
			if !ok {
				Err(ModuleRule, "对应Key", sk, "映射不存在，无法更新!")
				return FAILED
			}
			// 调用函数
			vr = nil
			if kon != UNDEFINED {
				Info(ModuleRule, ">>>>> 跳过带条件的赋值", sp.Setting.Ops)
				continue
			}
			if !IsConfig(p.Value[0], t) && !IsDigit(p.Value[0]) && IsCall(p.Value[0], t) {
				s := p.Value[0]
				RunBuiltin(vt, mid, key, s, r, t, ct, pc, p, vd)
				RunFuncs(vt, mid, key, s, r, t, ct, pc, p, vd)
				if SUCCESSFUL != RunBins(vt, mid, key, s, r, t, ct, pc, p, vd) {
					return FAILED
				}
			} else {
				vr = SetValue(mid, key, ct, vd, vp, p)
			}
			x, ok := vd.Vars[sk]
			// 赋值失败
			if !ok {
				Err(ModuleRule, "!!!!!!!!!", sk, "没有设置变量值")
				continue
			} else {
				Log(ModuleRule, "设置", sk, "值为", IfToString(x))
			}
			if vr == nil {
				vr = append(vr, x)
			}
			// 其他消息字段的赋值，所以要跳过
			if StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				sk := StrSearchStart(sk, AddDot(vd.Tops[key]...)...)
				sk = strings.TrimSuffix(sk, DOT)
				s := strings.Split(sk, DOT)
				kn := len(s)
				km, ok := vd.Msgs[sk]
				if !ok {
					continue
				}
				km = UpdateJsonOnKeys(km, vr, k[kn:])
				Info(ModuleRule, "更新消息字段", k, "的值为:", vr)
				if km != nil {
					vd.Msgs[sk] = km
					ts = []string{sk, RAFile}
					fk := strings.Join(ts, DOT)
					b, e := json.Marshal(km)
					if e == nil {
						vd.Vars[fk] = b
					}
				}
				continue
			}
			tm := UpdateJsonOnKeys(msg, vr, k[n:])
			Info(ModuleRule, "更新消息字段", k, "的值为:", vr)
			if tm != nil {
				msg = tm
				// 更新消息数据
				vd.Msgs[key] = msg
				b, e := json.Marshal(msg)
				if e == nil {
					fp.Content = b
					vd.Vars[fkey] = b
					vd.Vars[fkey2] = b
				}
			}
		}
	}
	// 更新消息数据
	vd.Msgs[key] = msg
	b, e := json.Marshal(msg)
	if e == nil {
		fp.Content = b
		vd.Vars[fkey] = b
		vd.Vars[fkey2] = b
	}
	return SUCCESSFUL
}

func updateDepValue(id, vt int, mid, key string, msg interface{}, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType,
	fp *FileType, vd *VarData, vp map[string]*VarType, sp *StepType) int {
	var vr []interface{}
	s := strings.Split(key, DOT)
	ts := []string{key, RAFile}
	fkey := strings.Join(ts, DOT)
	ts = []string{key, RAFileZ}
	fkey2 := strings.Join(ts, DOT)
	n := len(s)
	kon := SearchOp(OpOn, sp.Setting.Ops)
	// 更新赋值区域字段变量
	for i := range sp.Setting.Ops {
		// 跳过条件区
		p := &sp.Setting.Ops[i]
		if !FoundID(id, p.CIDs) {
			return SUCCESSFUL
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]

		if ok && tv.IsTemplate {
			continue
		}
		// 只赋值
		if !Equal(p.Op, OpEqual, OpConstruct) || len(p.Value) == 0 {
			continue
		}
		//其他消息变量字段
		if fp.Top == MainMsg && StrMatchStart(sk, vd.Tops[key]...) {
			x := StrSearchStart(sk, vd.Tops[key]...)
			_, ok := vd.Msgs[x]
			if !ok {
				Log(ModuleRule, "变量", sk, "消息未加载，跳过!")
				continue
			}
		}
		if StrMatchStart(sk, AddDot(key)...) {
			k, ok := vd.Keys[sk]
			if !ok {
				Err(ModuleRule, "对应Key", sk, "映射不存在，无法更新!")
				return FAILED
			}
			x, ok := vd.Vars[sk]
			if ok {
				if !HasKnownField(mid, key, t, p, vd) {
					Info(ModuleRule, "对应Key", sk, "已经赋值", x)
					continue
				}
			}
			// 调用函数
			vr = nil
			if kon != UNDEFINED {
				Info(ModuleRule, ">>>>> 跳过带条件的赋值", sp.Setting.Ops)
				continue
			}
			if !IsConfig(p.Value[0], t) && !IsDigit(p.Value[0]) && IsCall(p.Value[0], t) {
				s := p.Value[0]
				RunBuiltin(vt, mid, key, s, r, t, ct, pc, p, vd)
				RunFuncs(vt, mid, key, s, r, t, ct, pc, p, vd)
				if SUCCESSFUL != RunBins(vt, mid, key, s, r, t, ct, pc, p, vd) {
					return FAILED
				}
			} else {
				vr = SetValue(mid, key, ct, vd, vp, p)
			}
			x, ok = vd.Vars[sk]
			// 赋值失败
			if !ok {
				Err(ModuleRule, "!!!!!!!!!", sk, "没有设置变量值")
				continue
			} else {
				Log(ModuleRule, "设置", sk, "值为", IfToString(x))
			}
			if vr == nil {
				vr = append(vr, x)
			}
			// 其他消息字段的赋值，所以要跳过
			if StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				x := StrSearchStart(sk, AddDot(vd.Tops[key]...)...)
				x = strings.TrimSuffix(x, DOT)
				s := strings.Split(x, DOT)
				kn := len(s)
				km, ok := vd.Msgs[x]
				if !ok {
					continue
				}
				km = UpdateJsonOnKeys(km, vr, k[kn:])
				Info(ModuleRule, "更新消息字段", k, "的值为:", vr)
				if km != nil {
					vd.Msgs[x] = km
					ts = []string{x, RAFile}
					fk := strings.Join(ts, DOT)
					ts = []string{x, RAFileZ}
					fk2 := strings.Join(ts, DOT)
					b, e := json.Marshal(km)
					if e == nil {
						vd.Vars[fk] = b
						vd.Vars[fk2] = b
					}
				}
				continue
			}
			tm := UpdateJsonOnKeys(msg, vr, k[n:])
			Info(ModuleRule, "更新消息字段", k, "的值为:", vr)
			if tm != nil {
				msg = tm
				// 更新消息数据
				vd.Msgs[key] = msg
				b, e := json.Marshal(msg)
				if e == nil {
					fp.Content = b
					vd.Vars[fkey] = b
					vd.Vars[fkey2] = b
				}
			}
		}
	}
	// 更新校验区域字段变量
	for i := range sp.Expects {
		p := &sp.Expects[i]
		if !FoundID(id, p.CIDs) {
			continue
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]
		if ok && tv.IsTemplate {
			continue
		}
		for _, v := range p.Value {
			if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
				v = strings.Join([]string{mid, v}, DOT)
				LoadKeys(vd, vp, v)
			}
			// 其他消息变量字段
			if fp.Top == MainMsg && StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				continue
			}
			if Equal(v, fkey, fkey2) {
				Log(ModuleRule, "获取文件", fp.Name, "内容存到变量", v)
				vd.Vars[v] = fp.Content
				vd.Vars[fkey] = fp.Content
				vd.Vars[fkey2] = fp.Content
				continue
			}
			updateComposingKeys(id, vt, n, key, v, msg, vd)
		}
		// 其他消息变量字段
		if fp.Top == MainMsg && StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
			continue
		}
		updateComposingKeys(id, vt, n, key, sk, msg, vd)
	}
	// 更新消息数据
	vd.Msgs[key] = msg
	b, e := json.Marshal(msg)
	if e == nil {
		fp.Content = b
		vd.Vars[fkey] = b
		vd.Vars[fkey2] = b
	}
	return SUCCESSFUL
}

func updateMessage(id, vt int, mid, key string, msg interface{}, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType,
	fp *FileType, vd *VarData, vp map[string]*VarType, sp *StepType) int {
	var vr []interface{}
	s := strings.Split(key, DOT)
	ts := []string{key, RAFile}
	fkey := strings.Join(ts, DOT)
	ts = []string{key, RAFileZ}
	fkey2 := strings.Join(ts, DOT)
	extra := SearchOp(ExtraPresent, sp.Setting.Ops)
	Debug(ModuleRule, ">>>>>>> ++位置", extra)
	n := len(s)
	// 更新赋值区域字段变量
	start := SearchOp(OpOn, sp.Setting.Ops)
	prefix, m := SearchWildcardPaths(start, sp.Setting.Ops, vd)
	pm := map[string]interface{}{}
	if len(m) > 0 {
		if !WildcardConditionSatisfied(start, prefix, m, pm, sp.Setting.Ops, vp, vd) {
			Info(ModuleRunner, "不符合模糊校验条件，退出校验")
			return SUCCESSFUL
		}
	} else {
		if !ConditionSatisfied(sp.Setting.Ops, vp, vd) {
			Info(ModuleRunner, "不符合校验条件，退出校验")
			return SUCCESSFUL
		}
	}
	for i := range sp.Setting.Ops {
		// 跳过条件区
		if start != UNDEFINED && i >= start {
			break
		}
		p := &sp.Setting.Ops[i]
		if !FoundID(id, p.CIDs) {
			return SUCCESSFUL
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]
		// A -- x 删除字段要支持
		if ok && tv.IsTemplate && !Equal(p.Op, OpConstruct, NotPresent) {
			continue
		}
		// 只赋值
		if !Equal(p.Op, OpEqual, OpConstruct, NotPresent) || len(p.Value) == 0 {
			continue
		}
		// 用于依赖赋值判断，若变量值不存在，可以考虑赋值
		_, ok = vd.Vars[sk]
		// 其他消息变量字段
		if fp.Top == MainMsg && StrMatchStart(sk, AddDot(vd.Tops[key]...)...) && ok {
			continue
		}
		if StrMatchStart(sk, AddDot(key)...) || extra != UNDEFINED {
			k, ok := vd.Keys[sk]
			if !ok && extra == UNDEFINED {
				Err(ModuleRule, "对应Key", sk, "映射不存在，无法更新!")
				return FAILED
			}
			// 调用函数
			vr = nil
			if len(p.Value) > 0 && !IsConfig(p.Value[0], t) && !IsDigit(p.Value[0]) && IsCall(p.Value[0], t) {
				if start == UNDEFINED {
					continue
				}
				s := p.Value[0]
				RunBuiltin(vt, mid, key, s, r, t, ct, pc, p, vd)
				RunFuncs(vt, mid, key, s, r, t, ct, pc, p, vd)
				if SUCCESSFUL != RunBins(vt, mid, key, s, r, t, ct, pc, p, vd) {
					return FAILED
				}
			} else if p.Op != NotPresent && extra == UNDEFINED {
				if start == UNDEFINED {
					continue
				}
				vr = SetValue(mid, key, ct, vd, vp, p)
			}
			// 其他消息字段的赋值，所以要跳过
			if !StrMatchStart(sk, AddDot(key)...) || StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				continue
			}
			x, ok := vd.Vars[sk]
			// 赋值失败
			if !ok {
				Err(ModuleRule, "!!!!!!!!!", sk, "没有设置消息变量值")
				continue
			}
			if vr == nil {
				vr = append(vr, x)
			}
			tm := UpdateJsonOnKeys(msg, vr, k[n:])
			if tm != nil {
				msg = tm
			}
			if p.Op == NotPresent {
				Info(ModuleRule, ">>>>删除消息字段", p.Value)
				msg = DeleteMessageFields(sk, fkey, fkey2, key, pm, msg, ct, vd, vp, p, fp)
			} else if extra != UNDEFINED {
				if i > extra {
					kp := &sp.Setting.Ops[extra]
					Info(ModuleRule, ">>>>增加消息字段", sk, p.Value)
					msg = AddMessageFields(kp.Key, fkey, fkey2, key, pm, msg, ct, vd, vp, p, fp)
				}
			}
			// 更新消息数据
			b, e := json.Marshal(msg)
			if e == nil {
				fp.Content = b
				vd.Vars[fkey] = b
			}
		}
	}
	// 更新消息数据
	b, e := json.Marshal(msg)
	if e == nil {
		fp.Content = b
		vd.Vars[fkey] = b
	}
	return SUCCESSFUL
}

func UpdateRespVar(id, vt int, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType, ft *FlowType,
	tp *TemplateVarType, vd *VarData) int {
	// 需要加载变量对应的消息模板
	Log(ModuleRunner, "名称", tp.Key, "ID", tp.ID)
	s := []string{tp.ID, RAResp}
	h := strings.Join(s, DOT)
	b, ok := vd.Vars[h]
	if !ok {
		Err(ModuleRule, "!!!!!!无法获取到", tp.ID, "的应答!")
		return FAILED
	}
	msg := ToJSON(b.([]byte))
	if msg == nil {
		Log(ModuleRunner, "无法生成", h, "JSON数据!")
		return SUCCESSFUL
	}
	key := tp.ID
	// 存放应答消息，目前暂时不存，因为会洗掉之前存放的请求消息
	//vd.Msgs[key] = msg
	sn := strings.Split(tp.Key, DOT)
	sk := SearchLongestMatchNode(sn, t.Config.Nodes, DOT)
	service := t.Config.Nodes[sk]
	switch vt {
	case CaseVar:
		if sk != EMPTY {
			for i := range service.Steps {
				sp := &service.Steps[i]
				updateRespValue(id, vt, tp.ID, key, msg, r, t, ct, pc, vd, ct.Vars, sp)
			}
		}
		for i := range ct.Steps {
			sp := &ct.Steps[i]
			updateRespValue(id, vt, EMPTY, key, msg, r, t, ct, pc, vd, ct.Vars, sp)
		}
	case FlowVar:
		if sk != EMPTY {
			for i := range service.Steps {
				sp := &service.Steps[i]
				updateRespValue(id, vt, tp.ID, key, msg, r, t, ct, pc, vd, ct.Vars, sp)

			}
		}
		for i := range ft.Steps {
			sp := &ft.Steps[i]
			updateRespValue(id, vt, EMPTY, key, msg, r, t, ct, pc, vd, ct.Vars, sp)
		}
	case CFVar:
		for i := range pc.Flows {
			sp := &pc.Flows[i]
			updateRespValue(id, vt, EMPTY, key, msg, r, t, ct, pc, vd, pc.Vars, sp)
		}
	}

	return SUCCESSFUL
}

func updateRespValue(id, vt int, mid, key string, msg interface{},
	r TaskRequest, t *Testing, ct *CaseType, pc *ControlType,
	vd *VarData, vp map[string]*VarType, sp *StepType) int {
	s := strings.Split(key, DOT)
	n := len(s)
	for i := range sp.Setting.Ops {
		p := &sp.Setting.Ops[i]
		if !FoundID(id, p.CIDs) {
			continue
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]
		if ok && tv.IsTemplate {
			continue
		}
		// 更新value的值
		for _, v := range p.Value {
			if v == strings.TrimSpace(NoChars) {
				continue
			}
			if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
				v = strings.Join([]string{mid, v}, DOT)
				LoadKeys(vd, vp, v)
			}
			// 其他消息变量字段
			if StrMatchStart(v, vd.Tops[key]...) {
				continue
			}
			updateComposingKeys(id, vt, n, key, v, msg, vd)
		}
	}
	prefix := EMPTY
	for i := range sp.Expects {
		p := &sp.Expects[i]
		if !FoundID(id, p.CIDs) || p.Key == EMPTY {
			continue
		}
		sk := p.Key
		if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
			sk = strings.Join([]string{mid, p.Key}, DOT)
			LoadKeys(vd, vp, sk)
		}
		tv, ok := vp[sk]
		if ok && tv.IsTemplate {
			continue
		}
		// 获取模糊前缀
		if strings.Contains(sk, OpDotDotDot) && !StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
			s := strings.Split(sk, DOT)
			n := len(s)
			prefix = strings.Join(s[:n-1], DOT)
		}
		// 更新value的值
		for _, v := range p.Value {
			if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
				v = strings.Join([]string{mid, v}, DOT)
				LoadKeys(vd, vp, v)
			}
			// 其他消息变量字段
			if StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
				continue
			}
			updateComposingKeys(id, vt, n, key, v, msg, vd)

		}
		// 更新key的值
		if StrMatchStart(sk, AddDot(vd.Tops[key]...)...) {
			continue
		}
		updateComposingKeys(id, vt, n, key, sk, msg, vd)
		if !strings.Contains(sk, DOT) && prefix != EMPTY {
			s := []string{prefix, p.Key}
			v := strings.Join(s, DOT)
			k, ok := vd.Keys[v]
			if !ok {
				vd.Keys[v] = BuildKeys(v)
				k = vd.Keys[v]
			}
			updateWildcardKeys(id, vt, 1, v, msg, k, vd)
		}
	}
	return SUCCESSFUL
}
