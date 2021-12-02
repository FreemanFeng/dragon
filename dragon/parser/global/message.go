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

package global

import (
	"bytes"
	"regexp"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

/*
************************************
# 消息构造
************************************
// 构造列表类型，包含100个元素，每一个元素是字典类型，其中m.k.a字段的取值范围是[1,20]
> M1

x.y[100] -> m.k.a = 1..20 m.k.b = XTYPES m.s.c = FIELDS

// 构造列表类型，包含50个元素，每一个元素是字典类型，其中a字段的取值范围是[1,20]
> M2

m.k[50] -> a = 1..20 b = XTYPES c = kk

// 构造列表类型，包含所有可选值，每一项是字典类型
> M3

m.k[] -> a = 1..20 b = XTYPES c = kk

// 构造字典类型，字段值可以是范围数据
> M4

m.k -> a = 1..20 b = XTYPES c = kk

// 构造字典类型，字段值可来自构造数据，i.e. D20
> M5

m.k -> a = D20 b = XTYPES c = kk

// 构造列表类型，每一项是字典数据
> M6

[3] -> a = D20 b = XTYPES c = kk

// 构造列表类型，包含所有数据
> M7

[] -> a b c 1..9 30
*/

func ParseMessageConstructor(r TaskRequest, start int, content [][]byte, nodes map[string]*ConstructType) (int, int, string) {
	n := len(content)
	end := start
	code := SUCCESSFUL
	msg := MsgSuccessful
	rw := regexp.MustCompile(`\w+`)
	id := EMPTY
	for i := start; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		// 消息构造
		if MatchStart(c, MarkClass) {
			b := rw.Find(c)
			if len(b) == 0 {
				return end, FAILED, MsgIncompleteRule
			}
			id = string(b)
			_, ok := nodes[id]
			if ok {
				return end, FAILED, MsgDuplicatedName
			}
			Debug(ModuleParser, "Find Message Constructor ID", id)
			nodes[id] = &ConstructType{ID: id}
			continue
		}
		// 空行或注释
		if len(c) == 0 || MatchStart(c, MarkComment) || MatchStart(bytes.TrimSpace(c), MarkComment2) {
			continue
		}
		// 判断是否已获取消息ID
		r, ok := nodes[id]
		if id == EMPTY || !ok {
			return end, FAILED, MsgIncompleteRule
		}
		end = i
		k := bytes.Split(c, []byte(MarkConstruct))
		if len(k) < 2 {
			return end, FAILED, MsgIncompleteRule
		}
		r.Op = MarkConstruct
		c = bytes.TrimSpace(k[0])
		b := ConvertSpacesInsideQuotes(c)
		s := strings.ReplaceAll(string(b), BoundaryChars, SPACE)
		if StrMatchEnd(s, RightBracket) {
			r.IsList = true
		}
		c = bytes.TrimSpace(k[1])
		code, msg = ConstructMessage(s, c, r)
	}
	return end, code, msg
}
