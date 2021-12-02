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

package util

import (
	"testing"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func TestGetTimeFormat(t *testing.T) {
	s := GetTimeFormat(0, "2006")
	Log(s)
	s = GetTimeFormat(0, "06")
	Log(s)
	s = GetTimeFormat(0, "6")
	Log(s)
	s = GetTimeFormat(0, "006")
	Log(s)
	s = GetTimeFormat(0, "01")
	Log(s)
	s = GetTimeFormat(0, "02")
	Log(s)
	s = GetTimeFormat(0, "15")
	Log(s)
	s = GetTimeFormat(0, "04")
	Log(s)
	s = GetTimeFormat(0, "05")
	Log(s)
	s = GetTimeFormat(0, ".999")
	Log(s)
	s = GetTimeFormat(0, ".999999")
	Log(s)
	s = GetTimeFormat(0, ".999999999")
	Log(s)
	s = GetTimeFormat(0, "2006:01:02")
	Log(s)
	s = GetSeconds()
	Log(s)
	secs := CalculateSecs("30m")
	Log(secs)
	secs = CalculateSecs("30M")
	Log(secs)
	secs = CalculateSecs("30h")
	Log(secs)
	secs = CalculateSecs("30H")
	Log(secs)
	secs = CalculateSecs("30d")
	Log(secs)
	secs = CalculateSecs("30D")
	Log(secs)
}
