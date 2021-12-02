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
	"os/exec"
	"path"
	"reflect"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
	"github.com/FreemanFeng/dragon/dragon/util"
)

func buildParams(mid, key string, vp map[string]interface{}, t *Testing, op *OpType) []interface{} {
	var params []interface{}
	for pos, v := range op.Value {
		// 第一个是函数名，要bypass
		if pos == 0 {
			continue
		}
		if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
			v = strings.Join([]string{mid, v}, DOT)
		}
		k, ok := vp[v]
		if ok {
			//Log(ModuleRule, pos, ">>>>> 参数", v, "变量值", k, "类型", TypeOf(k))
			params = append(params, k)
		} else {
			//Log(ModuleRule, pos, ">>>>>!!! 参数", v, "类型", TypeOf(v))
			params = append(params, v)
		}
	}
	return params
}

func RunBuiltin(vt int, mid, vk, key string, r TaskRequest, t *Testing,
	ct *CaseType, pc *ControlType, op *OpType, vd *VarData) {
	if !IsBuiltin(key) || HasUnknownField(mid, vk, t, op, vd) {
		return
	}
	Log(ModuleRunner, "++++++++++++++++ 调用内置函数 +++++++++++++++++++")
	f := GetBuiltin(key)
	params := buildParams(mid, vk, vd.Vars, t, op)
	k := util.Call(f, params)
	if op.Key != EMPTY && len(k) > 0 {
		sk := op.Key
		if mid != EMPTY && !StrMatchStart(op.Key, CtxP) {
			sk = strings.Join([]string{mid, op.Key}, DOT)
		}
		Log(ModuleRunner, sk, "执行内置函数", key, "参数", params, "结果", k[0].Interface())
		vd.Vars[sk] = k[0].Interface()
	}
}

func IsCall(key string, t *Testing) bool {
	Debug(ModuleRunner, "判断", key, "是否函数调用")
	if IsReserved(key) || IsBuiltin(key) || IsFunc(key, t) || IsBin(key, t) {
		return true
	}
	return false
}

func IsConfig(key string, t *Testing) bool {
	Debug(ModuleRunner, "判断", key, "是否静态变量、消息构造变量、数据构造变量")
	_, ok1 := t.Config.Settings[key]
	_, ok2 := t.Config.Data[key]
	_, ok3 := t.Config.Messages[key]
	_, ok4 := t.Config.Meta[key]
	if ok1 || ok2 || ok3 || ok4 {
		return true
	}
	return false
}

func IsFunc(key string, t *Testing) bool {
	_, ok := t.Funcs[key]
	if !ok || IsReserved(key) {
		return false
	}
	return true
}

func IsBin(key string, t *Testing) bool {
	_, ok := t.Bins[key]
	if !ok || IsReserved(key) {
		return false
	}
	return true
}

func RunFuncs(vt int, mid, vk, key string, r TaskRequest, t *Testing,
	ct *CaseType, pc *ControlType, op *OpType, vd *VarData) {
	var k []reflect.Value
	p, ok := vd.Funcs[key]
	if !ok {
		p, ok = t.Funcs[key]
	}
	if !ok || IsReserved(key) || HasUnknownField(mid, vk, t, op, vd) {
		return
	}
	f, _ := p.Funcs[key]
	Log(ModuleRunner, "++++++++++++++++ 调用插件函数 +++++++++++++++++++")
	params := buildParams(mid, vk, vd.Vars, t, op)
	Log(ModuleRunner, "准备执行插件函数", key, "参数", params)
	if p.IsRPC {
		if IsGoCode(p) {
			k = util.CallFunc(f, p.Port, key, JsonListByte(params...))
		} else {
			k = util.CallFunc(f, p.Port, key, JsonListStr(params...))
		}
	} else {
		k = util.CallFunc(f, params...)
	}
	if op.Key != EMPTY && len(k) > 0 {
		sk := op.Key
		if mid != EMPTY && !StrMatchStart(op.Key, CtxP) {
			sk = strings.Join([]string{mid, op.Key}, DOT)
		}
		Log(ModuleRunner, sk, "执行插件函数", key, "参数", params, "结果", k[0].Interface())
		vd.Vars[sk] = k[0].Interface()
	}
}

func RunCallback(key string, t *Testing, vd *VarData, params ...interface{}) []reflect.Value {
	p, ok := vd.Funcs[key]
	if !ok {
		p, ok = t.Funcs[key]
	}
	if !ok {
		Err(ModuleRunner, ">>>>>>> 无法找到", key, "回调函数!")
		return nil
	}
	f, ok := p.Funcs[key]
	if p.IsRPC {
		if IsGoCode(p) {
			return util.CallFunc(f, p.Port, key, JsonListByte(params...))
		}
		// 其他语言无法识别JSON序列化[]byte
		return util.CallFunc(f, p.Port, key, JsonListStr(params...))
	} else {
		return util.CallFunc(f, params...)
	}
}

func HasCallback(key string, t *Testing, vd *VarData) bool {
	_, ok := vd.Funcs[key]
	return ok
}

func RunBins(vt int, mid, vk, key string, r TaskRequest, t *Testing,
	ct *CaseType, pc *ControlType, op *OpType, vd *VarData) int {
	bin, ok := vd.Bins[key]
	if !ok {
		bin, ok = t.Bins[key]
	}
	if !ok || IsReserved(key) || HasUnknownField(mid, vk, t, op, vd) {
		return SUCCESSFUL
	}
	Log(ModuleRunner, "++++++++++++++++ 调用可执行文件 +++++++++++++++++++")
	v := buildParams(mid, key, vd.Vars, t, op)
	x := path.Ext(bin)
	var params []string
	if x == ExtBash || x == ExtSh || x == ExtZsh {
		params = append(params, bin)
		bin = strings.Split(x, DOT)[1]
	}
	if x == ExtPy {
		params = append(params, bin)
		bin = CmdPython
	}
	if x == ExtJar {
		params = append(params, ArgJar, bin)
		bin = CmdJava
	}
	for _, k := range v {
		params = append(params, k.(string))
	}
	cmd := exec.Command(bin, params...)
	k, err := cmd.Output()
	if err != nil {
		Err(ModuleRunner, "执行", bin, "出错", err)
		return FAILED
	}
	if op.Key != EMPTY {
		sk := op.Key
		if mid != EMPTY && !StrMatchStart(op.Key, CtxP) {
			sk = strings.Join([]string{mid, op.Key}, DOT)
		}
		Log(ModuleRunner, sk, "执行可执行文件", key, "参数", params, "结果", string(k))
		vd.Vars[sk] = string(bytes.TrimSpace(k))
	}
	return SUCCESSFUL
}

func HasUnknownField(mid, key string, t *Testing, op *OpType, vd *VarData) bool {
	sk := op.Key
	if mid != EMPTY && !StrMatchStart(sk, CtxP) {
		sk = strings.Join([]string{mid, op.Key}, DOT)
	}
	for _, v := range op.Value {
		if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
			v = strings.Join([]string{mid, v}, DOT)
		}
		if strings.Contains(v, DOT) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
			_, ok := vd.Vars[v]
			if !ok {
				Err(ModuleRule, sk, "!!!!!!无法识别依赖变量", v, "key", key)
				return true
			}
		}
	}
	return false
}

func HasKnownField(mid, key string, t *Testing, op *OpType, vd *VarData) bool {
	sk := op.Key
	if mid != EMPTY && !StrMatchStart(sk, CtxP) {
		sk = strings.Join([]string{mid, op.Key}, DOT)
	}
	for _, v := range op.Value {
		if mid != EMPTY && !StrMatchStart(v, CtxP) && !IsConfig(v, t) && !IsDigit(v) && !IsCall(v, t) {
			v = strings.Join([]string{mid, v}, DOT)
		}
		if StrMatchStart(v, AddDot(vd.Tops[key]...)...) || StrMatchStart(v, AddDot(key)...) {
			Info(ModuleRule, sk, "存在依赖变量", v, "key", key)
			return true
		}
	}
	return false
}
