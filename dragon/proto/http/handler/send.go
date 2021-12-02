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

package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
	"github.com/nuxim/dragon/dragon/rule"
	"github.com/nuxim/dragon/dragon/util"
)

func OnSending(id int, t *Testing, p *VarData, params map[string]string, m, c, r map[string][]byte) int {
	var ms []map[string][]byte

	if !rule.HasCallback(HttpOnSending, t, p) {
		return Send(id, t, p, params, m, c, r)
	}

	x := rule.RunCallback(HttpOnSending, t, p, m, c, r)
	Info("回调", HttpOnSending, "函数")
	if len(x) > 0 {
		b := x[0].Interface().([]byte)
		ms = DecodeJsonList(b, p.Funcs[HttpOnSending])
	}
	n := len(ms)
	if n < 3 || ms == nil {
		Err(ModuleProto, "反序列化", HttpOnSending, "返回结果列表大小", n, "至少为3")
		return FAILED
	}

	for k, v := range ms[0] {
		m[k] = v
	}
	for k, v := range ms[1] {
		c[k] = v
	}
	for k, v := range ms[2] {
		r[k] = v
	}

	return SUCCESSFUL
}

func Send(id int, t *Testing, p *VarData, params map[string]string, m, c, r map[string][]byte) int {
	config := t.Config
	ka := config.Attributes[Args]
	ab := m[ka]
	kh := config.Attributes[Header]
	hb := m[kh]
	km := config.Attributes[MainMsg]
	bb := m[km]
	resp, code, head := util.RequestHttp(params[CFUrl], params[CFMethod], ab, hb, bb)
	r[RACode] = []byte(strconv.Itoa(code))
	for k, v := range head {
		s := []string{params[CFName], LowerCaseNoHyphen(k)}
		h := strings.Join(s, DOT)
		params[h] = v[0]
		s = []string{params[CFName], AddUnderScope(LowerCaseNoHyphen(k))}
		h = strings.Join(s, DOT)
		params[h] = v[0]
	}
	bh, err := json.Marshal(head)
	if err != nil {
		Err(ModuleProto, "序列化应答头", head, "出错", err)
		return FAILED
	}
	r[RAHead] = bh
	r[RAResp] = resp
	return SUCCESSFUL
}
