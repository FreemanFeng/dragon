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

	. "github.com/FreemanFeng/dragon/dragon/common"
)

/*
************************************
# 数据构造
************************************
> D1

x.y.z -> x = 1..100 y = 1..100 z = 1..100

> D2

xyz -> x = 1..100 y = 1..100 z = 1..100

> D3

"x y z" -> x = FIELDS y = TYPES z = XFIELDS

> D4

x:y.z -> x = FIELDS y = TYPES z = XFIELDS

> D5

"a=x b=x" -> x = FIELDS

> D6

x:y.z -> x = FIELDS y = TYPES z = XFIELDS

> D7

[10] = 1

> D8

[10] = DIGITS

> D9

[10] = 1..100

> D10

[10] x*3 x*7 -> x = 1..100

> D11

[5] x/3+y*3 y/7+x*2 -> x = 1..100 y = 1..5

> D12

[] = FIELDS

> D13

[10] = kkk

> D14

[] x_k -> x = FIELDS

> D15

[5] x_k -> x = FIELDS

> D16

[10] x:y.z -> x = FIELDS y = TYPES z = XFIELDS

> D17

[10] x:y.z -> x = FIELDS y = TYPES z = XFIELDS

> D18

[10]  =  1..100

> D19

[5] = 1..100 y = 1..5

> D20

x/9 -> x = 1..20

> D21

[] x/9 -> x = 1..20

> D22

[] x/9 (x*3+2)/5 -> x = 1..20
*/

func ParseDataConstructors(r TaskRequest, start int, content [][]byte, nodes map[string]*ConstructType) (int, int, string) {
	n := len(content)
	end := start
	code := SUCCESSFUL
	msg := MsgSuccessful
	rw := regexp.MustCompile(`\w+`)
	id := EMPTY
	rs := regexp.MustCompile(`\s+`)
	op := MarkConstruct
	for i := start; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		// 数据构造
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
			Debug(ModuleParser, "Find Data Constructor ID", id)
			nodes[id] = &ConstructType{ID: id, Formats: []FormatType{}}
			continue
		}
		// 空行或注释
		if len(c) == 0 || MatchStart(c, MarkComment) || MatchStart(bytes.TrimSpace(c), MarkComment2) {
			continue
		}
		// 判断是否已获取数据ID
		_, ok := nodes[id]
		if id == EMPTY || !ok {
			return end, FAILED, MsgIncompleteRule
		}
		end = i
		k := bytes.Split(c, []byte(MarkConstruct))
		if len(k) < 2 {
			k = bytes.Split(c, []byte(OpEqual))
			op = OpEqual
			if len(k) < 2 {
				return end, FAILED, MsgIncompleteRule
			}
		}
		c = bytes.TrimSpace(k[0])
		b := ConvertSpacesInsideQuotes(c)
		data := rs.Split(string(b), UNDEFINED)
		m := len(data)
		s := strings.ReplaceAll(data[0], BoundaryChars, SPACE)
		r := nodes[id]
		r.Op = op
		switch m {
		case 1:
			if strings.Index(s, LeftBracket) == 0 {
				r.IsList = true
			} else {
				f := BuildFormat(s)
				r.Formats = append(r.Formats, f)
			}
		default:
			if strings.Index(s, LeftBracket) != 0 {
				return end, FAILED, MsgUnknownRule
			}
			r.IsList = true
			for i := 1; i < m; i++ {
				k := strings.ReplaceAll(data[i], BoundaryChars, SPACE)
				f := BuildFormat(k)
				r.Formats = append(r.Formats, f)
			}
		}
		c = bytes.TrimSpace(k[1])
		code, msg = ConstructData(s, c, r)
	}
	return end, code, msg
}
