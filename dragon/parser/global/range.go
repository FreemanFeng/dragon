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
# 范围
************************************
XTYPES = haha kk xx

TYPES = patch hotfix release

FIELDS = repo version type

XFIELDS = nono coco

SPACES = "  "

MIXED = "A B " kk aa bb

DIGITS = 3 9..20 5
*/
func ParseRange(r TaskRequest, start int, content [][]byte, nodes map[string][]string) (int, int, string) {
	rs := regexp.MustCompile(`\s+`)
	n := len(content)
	end := start
	for i := start; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		if len(c) == 0 || MatchStart(c, MarkComment) || MatchStart(bytes.TrimSpace(c), MarkComment2) {
			continue
		}
		end = i
		// i.e. Types = patch release hotfix
		b := ConvertSpacesInsideQuotes(c)
		k := strings.Split(string(b), OpEqual)
		n := len(k)
		if n < 2 {
			return end, FAILED, MsgIncompleteRule
		}
		key := strings.TrimSpace(strings.ReplaceAll(k[0], BoundaryChars, SPACE))
		t := rs.Split(strings.TrimSpace(k[1]), UNDEFINED)
		// 相同变量名，值会累加
		nodes[key] = append(nodes[key], t...)
		n = len(nodes[key])
		for i := 0; i < n; i++ {
			s := strings.ReplaceAll(nodes[key][i], BoundaryChars, SPACE)
			nodes[key][i] = TrimDoubleQuotes(s)
		}
	}
	return end, SUCCESSFUL, MsgSuccessful
}
