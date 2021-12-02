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
	"bytes"
	"strings"
)

func MatchStart(c []byte, marks ...string) bool {
	for _, v := range marks {
		if bytes.Index(c, []byte(v)) == 0 {
			return true
		}
	}
	return false
}

func StrMatchStart(c string, marks ...string) bool {
	for _, v := range marks {
		if strings.Index(c, v) == 0 {
			return true
		}
	}
	return false
}

func AddDot(marks ...string) []string {
	var s []string
	for _, v := range marks {
		s = append(s, strings.Join([]string{v, DOT}, EMPTY))
	}
	return s
}

func StrMatch(c string, marks ...string) bool {
	for _, v := range marks {
		if strings.Index(c, v) >= 0 {
			return true
		}
	}
	return false
}

func StrSearchStart(c string, marks ...string) string {
	for _, v := range marks {
		if strings.Index(c, v) == 0 {
			return v
		}
	}
	return EMPTY
}

func MatchEnd(c []byte, marks ...string) bool {
	n := len(c)
	for _, v := range marks {
		if n-bytes.LastIndex(c, []byte(v)) == len([]byte(v)) {
			return true
		}
	}
	return false
}

func StrMatchEnd(c string, marks ...string) bool {
	n := len(c)
	for _, v := range marks {
		if n-strings.LastIndex(c, v) == len(v) {
			return true
		}
	}
	return false
}

func MatchInside(c []byte, marks ...string) bool {
	for _, v := range marks {
		if bytes.Index(c, []byte(v)) >= 0 {
			return true
		}
	}
	return false
}

func Find(c []byte, marks ...string) int {
	for _, v := range marks {
		x := bytes.Index(c, []byte(v))
		if x >= 0 {
			return x
		}
	}
	return UNDEFINED
}

func Equal(c string, marks ...string) bool {
	for _, v := range marks {
		if v == c {
			return true
		}
	}
	return false
}

func CombineOps() []string {
	return []string{OpOn, OpWithout, OpWith, NotPresent, ExtraPresent, OpConstruct, OpColon, OR}
}
func BasicOps() []string {
	return []string{OpEqual, MoreThan, LessThan, NotEqual}
}

func IsAction(key string) bool {
	if Equal(key, AKHttp, AKHttps, AKCache, ActKeyDB, AKReceive, AKSend, AKMock,
		ActKeyWEB, ActKeyCOM, ActKeyIOS, AKAndroid, ActKeyROS) {
		return true
	}
	return false
}

func HasKey(key string, s []string) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}
	return false
}

func IsHttpMethod(key string) bool {
	if Equal(key, MethodPost, MethodGet, MethodHead,
		MethodPut, MethodPatch, MethodDelete, MethodOptions) {
		return true
	}
	return false
}

func IsReserved(key string) bool {
	return Equal(key, REQUEST, REQUEST2, REQUEST3)
}

func IsBuiltin(key string) bool {
	return Equal(key, FuncNameFormatDate, FuncNameDATE, FuncNameFormatTime,
		FuncNameTIME, FuncNameFormatSeconds, FuncNameSEC, FuncNameFormatNowInSeconds,
		FuncNameNOW, FuncNameNowAddSeconds, FuncNameNowMS, FuncNameMD5, FuncNameSHA1, FuncNameBase64,
		FuncNameRandomString, FuncNameRandomDigit, FuncNameRandom, FuncNameUUID,
		FuncNameLog)
}

func GetBuiltin(key string) interface{} {
	return GetCall(Join(COLON, ModuleBuiltin, key))
}
