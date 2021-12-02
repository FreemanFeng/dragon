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
	"strings"
)

func IsPlugin(name string) bool {
	t := strings.Split(name, DOT)
	n := len(t)
	s := t[n-1]
	if s == ZIP || s == GZ || s == TGZ {
		return true
	}
	return false
}

func ZipType(name string) string {
	t := strings.Split(name, DOT)
	n := len(t)
	s := t[n-1]
	if s == GZ && n > 2 && t[n-2] == TAR {
		return TGZ
	}
	if s == ZIP || s == SZ || s == TGZ {
		return s
	}
	return EMPTY
}
