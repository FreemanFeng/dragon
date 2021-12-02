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

package http

import (
	"strings"

	"github.com/FreemanFeng/dragon/dragon/proto/http/handler"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func Run() {
	//日期时间
	SetCall(Join(COLON, ModuleProto, ProtoHTTP, FuncNameRequest), Request)
}

func Request(id, vid int, proto string, sp *NodeType, tr TaskRequest, t *Testing, ct *CaseType, pc *ControlType,
	op *OpType, vt *TemplateVarType, p *VarData) int {
	var hosts []string
	Info(ModuleProto, "并发id", id, "请求", vt.ID, vt.Key)
	if sp.Method == EMPTY || sp.Path == EMPTY {
		Err(ModuleProto, "无法识别http方法及路径")
		return FAILED
	}
	method := strings.ToUpper(sp.Method)
	path := sp.Path

	if len(sp.Hosts) > 0 {
		for _, v := range sp.Hosts {
			if k, ok := t.Config.Settings[v]; ok {
				hosts = append(hosts, k...)
			} else {
				hosts = append(hosts, v)
			}
		}
	} else {
		hosts = append(hosts, p.Host)
	}
	if !IsHttpMethod(method) {
		Err(ModuleProto, "Http方法", method, "不支持!")
		return FAILED
	}
	Log(ModuleRunner, "Http请求路径:", path, "请求方法:", method, "hosts:", hosts, "消息模板名:", vt.Key)
	mf := map[string]*FileType{}
	m := map[string][]byte{}
	c := map[string][]byte{}
	r := map[string][]byte{}
	params := map[string]string{}
	for _, v := range p.Data[vt.ID] {
		if v.Top == MainMsg && v.Name != vt.Key {
			continue
		}
		mf[v.Top] = v
		m[v.Top] = v.Content
		Info(ModuleProto, "id:", id, "请求", v.Name, "top:", v.Top, "path:", v.Path)
	}
	for _, host := range hosts {
		s := []string{proto, "://", host, path}
		reqUrl := strings.Join(s, EMPTY)
		Info(ModuleProto, ">>>>>>> 请求URL:", reqUrl)
		params[CFName] = vt.ID
		params[CFUrl] = reqUrl
		params[CFMethod] = method
		s = []string{params[CFName], RAResp}
		params[RAResp] = strings.Join(s, DOT)
		s = []string{params[CFName], RACode}
		params[RACode] = strings.Join(s, DOT)
		s = []string{params[CFName], RAHead}
		params[RAHead] = strings.Join(s, DOT)
		if SUCCESSFUL != handler.OnReadySending(id, t, p, params, m, c, r) {
			return FAILED
		}
		if SUCCESSFUL != handler.OnSending(id, t, p, params, m, c, r) {
			return FAILED
		}
		if SUCCESSFUL != handler.OnReceived(id, t, p, params, m, c, r) {
			return FAILED
		}
		if SUCCESSFUL != handler.OnError(id, t, p, params, m, c, r) {
			return FAILED
		}
		bn := params[RAResp]
		p.Vars[bn] = r[RAResp]
		cn := params[RACode]
		p.Vars[cn] = r[RACode]
		hn := params[RAHead]
		p.Vars[hn] = r[RAHead]
		// 应答头
		for k, v := range params {
			if strings.Contains(k, DOT) {
				p.Vars[k] = v
			}
		}
		Info(ModuleProto, ">>>>>> id:", id, "应答内容", bn, ":", string(r[RAResp]),
			"应答状态码", cn, ":", string(r[RACode]))
	}
	return SUCCESSFUL
}
