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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGobEncode(t *testing.T) {
	r := assert.New(t)
	m := LogConfig{LogLevel: INFO, LogALL: 1, LogModules: map[string]int{ModuleDB: 1, ModuleCommon: 1}}
	p := GobEncode(m)
	var q LogConfig
	GobDecode(p, &q)
	Log(q)
	r.Equal(m, q, EMPTY)
}
