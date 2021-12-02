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
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
	"github.com/nuxim/dragon/dragon/proto/http"
	"github.com/nuxim/dragon/dragon/rule"
)

func Request(id int, r TaskRequest, t *Testing, ct *CaseType, pc *ControlType, ft *FlowType,
	op *OpType, tp *TemplateVarType, cd, pd, fd *VarData) int {
	Log(ModuleRunner, ">>>>>>>>>>>> 并发id", id, "请求", tp.ID, tp.Key, "ft:", ft)
	// 可以通过文件名指定协议
	proto := ct.Proto
	k, ok := t.Templates[tp.Key]
	if ok {
		if HasKey(AKHttp, k.Actions) {
			proto = HTTP
		}
		if HasKey(AKHttps, k.Actions) {
			proto = HTTPS
		}
	}
	p := fd
	vt := FlowVar
	if ft == nil {
		p, vt = selectData(tp, cd, pd)
	}
	if p == nil {
		Info(ModuleRunner, "无法识别变量", tp.ID)
		return FAILED
	}
	if SUCCESSFUL != rule.UpdateVar(id, vt, proto, r, t, ct, pc, ft, tp, p) {
		return FAILED
	}
	// 上下文模板无需发请求
	if tp.IsCtx {
		return SUCCESSFUL
	}
	// i.e. GetWeather.json 获取GetWeather对应的服务配置
	s := strings.Split(tp.Key, DOT)
	key := SearchLongestMatchNode(s, t.Config.Nodes, DOT)
	if key == EMPTY {
		// fallback逻辑
		key = SearchLongestMatchSetting(s, t.Config.Settings, DOT)
		if key == EMPTY || len(t.Config.Settings[key]) < 2 {
			Err(ModuleRunner, tp.Key, "没有对应的节点定义!")
			return FAILED
		}
		v := t.Config.Settings[key]
		method := v[0]
		path := v[1]
		t.Config.Nodes[key] = &NodeType{NID: key, Proto: proto, Path: path, Method: method}
	}
	service := t.Config.Nodes[key]
	Log(ModuleRunner, tp.ID, "服务key", key)
	switch proto {
	case HTTP:
		if SUCCESSFUL != http.Request(id, vt, proto, service, r, t, ct, pc, op, tp, p) {
			return FAILED
		}
	}
	if SUCCESSFUL != rule.UpdateRespVar(id, vt, r, t, ct, pc, ft, tp, p) {
		return FAILED
	}
	return SUCCESSFUL
}

func selectData(tp *TemplateVarType, cd, pd *VarData) (*VarData, int) {
	// 默认是控制流变量
	p := pd
	vt := CFVar
	_, ok := pd.Data[tp.ID]
	if !ok {
		// 判断是否用例变量
		_, ok = cd.Data[tp.ID]
		if !ok {
			Log(ModuleRunner, "无法识别", tp.ID, "并生成相应请求!")
			return nil, UNDEFINED
		}
		p = cd
		vt = CaseVar
	}
	return p, vt
}
