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
	"regexp"
	"strconv"
	"strings"
)

func ParseSetting(c []byte) (int, string, SettingType) {
	if len(c) == 0 {
		return SUCCESSFUL, MsgSuccessful, SettingType{}
	}
	setting := SettingType{}
	code, msg, ops := ParseOp(c)
	if len(ops) == 0 {
		return code, msg, setting
	}
	setting.Ops = append(setting.Ops, ops...)
	Log(ModuleParser, ParserTC, "用例设置", setting.Ops)
	return code, msg, setting
}

func ParseOp(c []byte) (int, string, []OpType) {
	var ops []OpType
	rs := regexp.MustCompile(`\s+`)
	t := TrimSpaces(c)
	Log(ModuleParser, string(t))
	t = rs.ReplaceAll(t, []byte(SPACE))
	//m := map[int]bool{}
	if len(t) == 0 {
		return FAILED, MsgIncompleteRule, ops
	}
	k := bytes.Split(t, []byte(SPACE))
	for _, v := range k {
		ops = SplitOps(v, ops)
	}
	for i := range ops {
		v := ops[i].Value
		n := len(v)
		if n == 1 && len(v[0]) == 2 && StrMatchStart(v[0], LeftBrace) && StrMatchEnd(v[0], RightBrace) {
			ops[i].Value = nil
			ops[i].IsMap = true
		}
		if n < 2 {
			continue
		}
		if v[0] == LeftBracket && v[n-1] == RightBracket {
			ops[i].Value = v[1 : n-1]
			ops[i].IsList = true
			//Log("替换后", ops[i].Value)
		}
	}
	// 对于with/without操作符，若不是排序，则合并前一个操作
	// i.e. A...message with 暴雨 on city = 广州 district = 白云
	found := false
	for i := range ops {
		if i == 0 {
			continue
		}
		prev := &ops[i-1]
		p := &ops[i]
		if prev.Key == EMPTY && prev.Op == EMPTY && p.Key == EMPTY &&
			len(prev.Value) == 1 && Equal(p.Op, OpWith, OpWithout) {
			p.Key = prev.Value[0]
			prev.Value = nil
			found = true
		}
	}
	if found {
		var ts []OpType
		for i := range ops {
			p := &ops[i]
			if p.Key == EMPTY && p.Op == EMPTY && p.Value == nil {
				continue
			}
			ts = append(ts, ops[i])
		}
		return SUCCESSFUL, MsgSuccessful, ts
	}

	return SUCCESSFUL, MsgSuccessful, ops
}

func ConstructData(s string, c []byte, r *ConstructType) (int, string) {
	code := SUCCESSFUL
	msg := MsgSuccessful
	code, msg, r.Settings = ParseOp(c)
	for i := range r.Settings {
		TrimOp(&r.Settings[i])
	}
	if !r.IsList {
		Debug(ModuleParser, r.Formats, r.Settings)
		return code, msg
	}
	h := strings.Split(s, DOT)
	for _, v := range h {
		t := strings.Index(v, LeftBracket)
		n := strings.Index(v, RightBracket)
		if n == UNDEFINED {
			return FAILED, MsgIncompleteRule
		}
		if n == t+1 {
			r.Keys = append(r.Keys, BuildFormat(ZERO))
		} else {
			r.Keys = append(r.Keys, BuildFormat(v[t+1:n]))
		}
	}
	Debug(ModuleParser, r.Formats, r.Keys, r.Settings)
	return SUCCESSFUL, MsgSuccessful
}

func BuildKeys(s string) []FormatType {
	var keys []FormatType
	h := strings.Split(s, DOT)
	for _, v := range h {
		t := strings.Index(v, LeftBracket)
		n := strings.Index(v, RightBracket)
		// 找不到返回-1，仅当其中一个返回-1，乘积小于0，或者n==0
		if t*n < 0 || n == 0 {
			return keys
		}
		if t > 0 {
			keys = append(keys, BuildFormat(v[:t]))
			keys = append(keys, BuildFormat(v[t+1:n]))
			continue
		} else if t == 0 && n == 1 {
			v = ZERO
		} else if t == 0 && n > 1 {
			v = v[t+1 : n]
			_, e := strconv.Atoi(v)
			// 不可识别，默认为0
			if e != nil {
				v = ZERO
			}
		}
		keys = append(keys, BuildFormat(v))
	}
	return keys
}

func ConstructMessage(s string, c []byte, r *ConstructType) (int, string) {
	code := SUCCESSFUL
	msg := MsgSuccessful
	r.Keys = BuildKeys(s)
	if code != SUCCESSFUL {
		return code, msg
	}
	code, msg, r.Settings = ParseOp(c)
	for i := range r.Settings {
		TrimOp(&r.Settings[i])
	}
	Debug(ModuleParser, r.Keys, r.Settings)
	return code, msg
}
