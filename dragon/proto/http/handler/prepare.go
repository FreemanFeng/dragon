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
	"github.com/FreemanFeng/dragon/dragon/rule"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func OnReadySending(id int, t *Testing, p *VarData, params map[string]string, m, c, r map[string][]byte) int {
	var ms []map[string][]byte
	if !rule.HasCallback(HttpOnSending, t, p) {
		Info("没有", HttpOnReadySending, "回调函数")
		return SUCCESSFUL
	}
	x := rule.RunCallback(HttpOnReadySending, t, p, m, c, r)
	Info("回调", HttpOnReadySending, "函数")
	if len(x) > 0 {
		b := x[0].Interface().([]byte)
		ms = DecodeJsonList(b, p.Funcs[HttpOnReadySending])
	}
	n := len(ms)
	if n < 3 || ms == nil {
		Err(ModuleProto, "反序列化", HttpOnReadySending, "返回结果列表大小", n, "至少为3")
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
