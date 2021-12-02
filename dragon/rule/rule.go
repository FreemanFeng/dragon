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
	"encoding/json"
	"math/rand"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func ConditionSatisfied(ops []OpType, vt map[string]*VarType, vd *VarData) bool {
	var exp string
	var rev interface{}
	start := SearchOp(OpOn, ops)
	kor := SearchOp(OR, ops)
	if start == UNDEFINED || len(ops) == 0 {
		return true
	}
	Debug(ModuleRule, ">>>>>> or位置", kor, "on位置", start)
	// 构造变量存放在VarType，变量字段值存放在VarData
	n := len(ops)
	passed := true
	for i := start; i < n; i++ {
		p := &ops[i]
		if p.Op == OR {
			if passed {
				return true
			}
			passed = true
			continue
		}
		for _, v := range p.Value {
			// 取右边值，接收值
			exp = v
			// 赋值变量
			k, ok := vd.Vars[v]
			if ok {
				exp = IfToString(k)
			} else {
				// 构造变量
				x, ok := vt[v]
				pos, _ := vd.VPos[v]
				if ok && len(x.Value) > pos {
					exp = IfToString(x.Value[pos])
				}
			}
			// 取左边值，期望值
			k, ok = vd.Vars[p.Key]
			if !ok {
				x, ok := vt[v]
				pos, _ := vd.VPos[v]
				if ok && len(x.Value) > pos {
					rev = x.Value[pos]
				} else {
					if kor == UNDEFINED {
						return false
					}
					passed = false
				}
			} else {
				rev = k
			}
			// 比较两边的值，仅字符串对比
			if !CheckExpected(exp, p.Op, rev) {
				if kor == UNDEFINED {
					return false
				}
				passed = false
			}
		}
	}
	return passed
}

func VerificationPassed(start int, ops []OpType, vt map[string]*VarType, vd *VarData, fds map[string]*VarData) bool {
	var rev, exp interface{}
	if len(ops) == 0 {
		return true
	}
	prefix, m := SearchWildcardPaths(start, ops, vd)
	tm := map[string]interface{}{}
	//Log("+++++++++++++ 1")
	if len(m) > 0 {
		if !WildcardConditionSatisfied(start, prefix, m, tm, ops, vt, vd) {
			Info(ModuleRunner, "不符合模糊校验条件，退出校验")
			return true
		}
	} else {
		if !ConditionSatisfied(ops, vt, vd) {
			Info(ModuleRunner, "不符合校验条件，退出校验")
			return true
		}
	}
	//Log("+++++++++++++ 2")
	// 构造变量存放在VarType，变量字段值存放在VarData
	for i := 0; i < start; i++ {
		p := &ops[i]
		for _, v := range p.Value {
			//Log("+++++++++++++ 3", v)
			// 取右边值，期望值
			// 赋值变量
			k, ok := vd.Vars[v]
			//Log("+++++++++++++ 3", v, k)
			if ok {
				if x, ok := k.(map[string]interface{}); ok {
					//Log("+++++++++++++ 3.1", v, x)
					exp = x
				} else if x, ok := k.([]interface{}); ok {
					//Log("+++++++++++++ 3.2", v, x)
					exp = x
				} else {
					//Log("+++++++++++++ 3.3", v, k)
					exp = k
				}
			} else if strings.Contains(v, DOT) {
				// 构造变量值
				x, ok := vt[v]
				pos, _ := vd.VPos[v]
				if ok && len(x.Value) > pos {
					exp = x.Value[pos]
				}
				//Log("+++++++++++++ 3.4", v, exp)
			}
			// 取左边值，实际值
			rev = p.Key
			k, ok = vd.Vars[p.Key]
			//Log("+++++++++++++ 4", p.Key, k)
			// 普通的字段比较
			if ok {
				if x, ok := k.(map[string]interface{}); ok {
					//Log("+++++++++++++ 4.1", p.Key, x)
					rev = x
				} else if x, ok := k.([]interface{}); ok {
					//Log("+++++++++++++ 4.2", p.Key, x)
					rev = x
				} else {
					//Log("+++++++++++++ 4.3", p.Key, k)
					rev = k
				}
			}
			if !strings.Contains(p.Key, OpDotDotDot) {
				//Log("+++++++++++++ 5", p.Key, k, ok)
				if !ok {
					x, ok := vt[p.Key]
					pos, _ := vd.VPos[p.Key]
					if ok && len(x.Value) > pos {
						rev = x.Value[pos]
					}
					//Log("+++++++++++++ 5.1", p.Key, rev)
				}
				//Log("+++++++++++++ 6", exp, p.Op, rev)
				if !CheckExpected(exp, p.Op, rev) {
					Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", p.Op, "期望", exp, "不匹配")
					return false
				}
				continue
			}
			found := false
			// 两边都是字典
			if exp == nil || rev == nil {
				Err(ModuleRunner, "期望", exp, "操作符", p.Op, "实际", rev, "校验失败!")
				return false
			}
			if IsMap(exp) {
				if CheckExpected(exp, p.Op, rev) {
					Info(ModuleRunner, ">>>>>>>> 模糊校验，实际", rev, "操作符", p.Op, "期望", exp, "匹配成功")
					found = true
				}
			} else if IsList(exp) {
				if CheckExpected(exp, p.Op, rev) {
					Info(ModuleRunner, ">>>>>>>> 模糊校验，实际", rev, "操作符", p.Op, "期望", exp, "匹配成功")
					found = true
				}
			} else {
				for _, tv := range tm {
					rev = tv
					if CheckExpected(IfToString(exp), p.Op, rev) {
						Info(ModuleRunner, ">>>>>>>> 模糊校验，实际", rev, "操作符", p.Op, "期望", exp, "匹配成功")
						found = true
					}
				}
			}
			if !found {
				Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", p.Op, "期望", exp, "不匹配")
				return false
			}
		}
	}
	return true
}

func SearchWildcardPaths(start int, ops []OpType, vd *VarData) (string, map[string]interface{}) {
	m := map[string]interface{}{}
	prefix := EMPTY
	for i := 0; i < start; i++ {
		p := &ops[i]
		if !strings.Contains(p.Key, OpDotDotDot) {
			continue
		}
		k, ok := vd.Vars[p.Key]
		if !ok {
			continue
		}
		s := strings.Split(p.Key, DOT)
		n := len(s)
		prefix = strings.Join(s[:n-1], DOT)
		x, ok := k.(map[string]interface{})
		if !ok {
			continue
		}
		for xk, xv := range x {
			m[xk] = xv
		}
		if len(m) > 0 {
			break
		}
	}
	return prefix, m
}

func WildcardConditionSatisfied(start int, prefix string, m, tm map[string]interface{}, ops []OpType, vt map[string]*VarType, vd *VarData) bool {
	var rev, exp interface{}
	if start == UNDEFINED || len(ops) == 0 || len(m) == 0 {
		return true
	}
	// 构造变量存放在VarType，变量字段值存放在VarData
	n := len(ops)
	top := 0
	kor := SearchOp(OR, ops)
	passed := true
	cm := map[string]int{}
	for i := start; i < n; i++ {
		op := &ops[i]
		if op.Op == OR {
			if passed {
				break
			}
			// 新的一组
			top = 0
			cm = map[string]int{}
			passed = true
			continue
		}
		if len(op.Value) == 0 {
			continue
		}
		for _, v := range op.Value {
			// 取右边值，期望值
			exp = v
			// 赋值变量
			k, ok := vd.Vars[v]
			if ok {
				exp = k
			} else if strings.Contains(v, DOT) {
				// 变量值
				x, ok := vt[v]
				pos, _ := vd.VPos[v]
				if ok && len(x.Value) > pos {
					exp = x.Value[pos]
				}
			}
			// 取左边值，实际值
			rev = op.Key
			k, ok = vd.Vars[op.Key]
			// 普通的字段比较
			if strings.Contains(op.Key, DOT) {
				if ok {
					rev = k
				} else if !ok {
					x, ok := vt[op.Key]
					pos, _ := vd.VPos[op.Key]
					if ok && len(x.Value) > pos {
						rev = x.Value[pos]
					}
				}
				if !CheckExpected(IfToString(exp), op.Op, rev) {
					if kor == UNDEFINED {
						Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", op.Op, "期望", exp, "不匹配")
						return false
					}
					passed = false
				}
				continue
			}
			// 拼接前缀
			s := []string{prefix, op.Key}
			h := strings.Join(s, DOT)
			k, ok = vd.Vars[h]
			if !ok {
				rev = op.Key
				if !CheckExpected(IfToString(exp), op.Op, rev) {
					if kor == UNDEFINED {
						Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", op.Op, "期望", exp, "不匹配")
						return false
					}
					passed = false
				}
				continue
			}
			x, ok := k.(map[string]interface{})
			// 单一值
			if !ok {
				if z, ok := k.(interface{}); ok {
					rev = z
				} else {
					rev = op.Key
				}
				if !CheckExpected(IfToString(exp), op.Op, rev) {
					if kor == UNDEFINED {
						Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", op.Op, "期望", exp, "不匹配")
						return false
					}
					passed = false
				}
				continue
			}
			top++
			found := false
			for xk, xv := range x {
				rev = xv
				if CheckExpected(IfToString(exp), op.Op, rev) {
					Info(ModuleRunner, ">>>>>>>> 模糊校验，实际", rev, "操作符", op.Op, "期望", exp, "匹配成功")
					_, ok := cm[xk]
					if !ok {
						cm[xk] = 1
					} else {
						cm[xk] += 1
					}
					found = true
				}
			}
			if !found {
				if kor == UNDEFINED {
					Err(ModuleRunner, "!!!!!!!!模糊校验，实际", rev, "操作符", op.Op, "期望", exp, "不匹配")
					return false
				}
				passed = false
			}
		}
	}
	if !passed {
		return false
	}
	// 没有模糊字段
	if top == 0 {
		for k, v := range m {
			tm[k] = v
		}
	} else {
		for k, v := range cm {
			if v != top {
				continue
			}
			_, ok := m[k]
			if ok {
				tm[k] = m[k]
			}
		}
	}
	return true
}

func FoundKey(key string, s []string) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}
	return false
}

func FoundID(id int, d []int) bool {
	if len(d) == 0 {
		return true
	}
	for _, v := range d {
		if v == id {
			return true
		}
	}
	return false
}

func SearchOp(key string, ops []OpType) int {
	if len(ops) == 0 {
		return UNDEFINED
	}
	for i, op := range ops {
		if op.Op == key {
			return i
		}
	}
	return UNDEFINED
}

func SetValue(mid, key string, ct *CaseType, vd *VarData, vp map[string]*VarType, p *OpType) []interface{} {
	var vr []interface{}
	n := len(p.Value)
	if n == 0 {
		return nil
	}
	i := rand.Intn(n)
	sk := p.Key
	if mid != EMPTY && !StrMatchStart(p.Key, CtxP) {
		sk = strings.Join([]string{mid, p.Key}, DOT)
	}
	if ct.Fill == FillAll {
		_, ok := vd.VPos[sk]
		if !ok {
			vd.VPos[sk] = UNDEFINED
		}
		vd.VPos[sk] = (vd.VPos[sk] + 1) % n
		i = vd.VPos[sk]
	}
	v := p.Value[i]
	if v == strings.TrimSpace(NoChars) {
		vd.Vars[sk] = EMPTY
		vr = append(vr, EMPTY)
		return vr
	}
	if IsOpVars(v) {
		x := EvalVarData(sk, v, vd)
		if v != x {
			vr = append(vr, x)
			return vr
		}
	}
	vt, ok := vp[v]
	if ok && vt.IsOpt {
		n = len(vt.Value)
		// 构造变量，遍历取值
		if vt.Op == OpConstruct || ct.Fill == FillAll {
			_, ok := vd.VPos[v]
			if !ok {
				vd.VPos[v] = UNDEFINED
			}
			vd.VPos[v] = (vd.VPos[v] + 1) % n
			i = vd.VPos[v]
			vd.Vars[sk] = vt.Value[i]
			if vt.Op == OpConstruct {
				return vt.Value
			}
			vr = append(vr, vt.Value[i])
			return vr
		}
		// 随机赋值
		i = rand.Intn(n)
		vd.Vars[sk] = vt.Value[i]
		vr = append(vr, vt.Value[i])
		return vr
	}
	// 函数调用返回值
	k, ok := vd.Vars[v]
	if ok {
		vd.Vars[sk] = k
		vr = append(vr, k)
	} else {
		vd.Vars[sk] = v
		vr = append(vr, v)
	}
	return vr
}

func DeleteMessageFields(prefix, fkey, fkey2, key string, pm map[string]interface{}, msg interface{}, ct *CaseType, vd *VarData,
	vp map[string]*VarType, p *OpType, fp *FileType) interface{} {
	n := len(p.Value)
	if n == 0 {
		return nil
	}
	s := []string{prefix, NotPresent}
	k := strings.Join(s, EMPTY)
	Info(ModuleRule, ">>>>>>准备删除 1，key", k)
	_, ok := vd.VPos[k]
	if !ok {
		vd.VPos[k] = UNDEFINED
	}
	vd.VPos[k] = (vd.VPos[k] + 1) % n
	i := vd.VPos[k]
	v := p.Value[i]
	// 模糊路径
	if len(pm) > 0 {
		Info(ModuleRule, ">>>>>>准备删除 2，key", k, pm)
		for kk := range pm {
			Info(ModuleRule, ">>>>>>准备删除 3，key", k, kk)
			s := []string{kk, v}
			h := strings.Join(s, DOT)
			Info(ModuleRule, ">>>>>>准备删除 4，h", h)
			keys := BuildKeys(h)
			s = strings.Split(key, DOT)
			kn := len(s)
			tm := DeleteJSON(msg, keys[kn:])
			Info(ModuleRule, ">>>>>>准备删除 5，keys", keys[kn:])
			if tm != nil {
				Info(ModuleRule, ">>>>>>准备删除 6，tm", tm)
				msg = tm
				// 更新消息数据
				b, e := json.Marshal(msg)
				if e == nil {
					Info(ModuleRule, ">>>>>>准备删除 7，msg", msg)
					fp.Content = b
					vd.Vars[fkey] = b
					vd.Vars[fkey2] = b
				}
			}
		}
		return msg
	}
	s = []string{prefix, v}
	h := strings.Join(s, DOT)
	Info(ModuleRule, ">>>>>>准备删除 4，h", h)
	keys := BuildKeys(h)
	s = strings.Split(key, DOT)
	kn := len(s)
	Info(ModuleRule, ">>>>>>准备删除 5，keys", keys[kn:])
	tm := DeleteJSON(msg, keys[kn:])
	if tm != nil {
		Info(ModuleRule, ">>>>>>准备删除 6，tm", tm)
		msg = tm
		// 更新消息数据
		b, e := json.Marshal(msg)
		if e == nil {
			Info(ModuleRule, ">>>>>>准备删除 7，msg", msg)
			fp.Content = b
			vd.Vars[fkey] = b
			vd.Vars[fkey2] = b
		}
	}
	return msg
}

func AddMessageFields(prefix, fkey, fkey2, key string, pm map[string]interface{}, msg interface{}, ct *CaseType, vd *VarData,
	vp map[string]*VarType, p *OpType, fp *FileType) interface{} {
	var kv []interface{}
	n := len(p.Value)
	if n == 0 {
		return nil
	}
	s := []string{prefix, ExtraPresent}
	k := strings.Join(s, EMPTY)
	Info(ModuleRule, ">>>>>>准备增加 1，key", k)
	_, ok := vd.VPos[k]
	if !ok {
		vd.VPos[k] = UNDEFINED
	}
	vd.VPos[k] = (vd.VPos[k] + 1) % n
	i := vd.VPos[k]
	for x := i; x < i+n; x++ {
		ki := x % n
		kv = append(kv, ConvertValue(p.Value[ki]))
	}
	Info(ModuleRule, ">>>>>>准备增加 2，kv", kv)
	// 模糊路径
	if len(pm) > 0 {
		for kk := range pm {
			s := []string{kk, p.Key}
			h := strings.Join(s, DOT)
			keys := BuildKeys(h)
			s = strings.Split(key, DOT)
			kn := len(s)
			Info(ModuleRule, ">>>>>>准备增加 3，keys", keys[kn:])
			tm := AddJSON(msg, kv, keys[kn:])
			if tm != nil {
				Info(ModuleRule, ">>>>>>准备增加 4，tm", tm)
				msg = tm
				// 更新消息数据
				b, e := json.Marshal(msg)
				if e == nil {
					Info(ModuleRule, ">>>>>>准备增加 5，msg", msg)
					fp.Content = b
					vd.Vars[fkey] = b
					vd.Vars[fkey2] = b
				}
			}
		}
		return msg
	}
	s = []string{prefix, p.Key}
	h := strings.Join(s, DOT)
	keys := BuildKeys(h)
	s = strings.Split(key, DOT)
	kn := len(s)
	Info(ModuleRule, ">>>>>>准备增加 6，key", h, "keys", keys[kn:])
	tm := AddJSON(msg, kv, keys[kn:])
	if tm != nil {
		Info(ModuleRule, ">>>>>>准备增加 7，tm", tm)
		msg = tm
		// 更新消息数据
		b, e := json.Marshal(msg)
		if e == nil {
			Info(ModuleRule, ">>>>>>准备增加 8，msg", msg)
			fp.Content = b
			vd.Vars[fkey] = b
			vd.Vars[fkey2] = b
		}
	}
	return msg
}
