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
package suites

import (
	"bytes"

	. "github.com/nuxim/dragon/dragon/common"
)

func ParseExpect(c []byte) (int, string, []byte, []OpType) {
	k := bytes.Split(c, []byte(MarkExpect))
	if len(k) < 2 || len(bytes.TrimSpace(k[1])) == 0 {
		return FAILED, MsgIncompleteRule, c, []OpType{}
	}
	code, msg, ops := ParseOp(k[1])
	Debug(ModuleParser, ParserTC, "用例检查点", ops)
	return code, msg, bytes.TrimSpace(k[0]), ops
}
