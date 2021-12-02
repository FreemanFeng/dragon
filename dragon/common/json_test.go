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
)

func TestIfToString(t *testing.T) {
	var a int8
	a = 8
	s := IfToString(a)
	Log("a", s)
}

func TestWalkJSON(t *testing.T) {
	s := `   { "m" : { "k": [ { "a" : 1, "b" : "c", "x": [1, 2, 3] }, { "a" : 2, "b" : "d", "x":[3,4,5] } ], "s": "hello" }, "w": "world"}`
	msg, fields := WalkJSON([]byte(s), "R")
	for k, v := range fields {
		Log("key", k, "value", v)
	}
	Log(msg)
	//keys := BuildKeys("m...a")
	keys := BuildKeys("m...x[0]")
	//var keys []FormatType
	//keys = append(keys, FormatType{Content: "m"})
	//keys = append(keys, FormatType{Content: "k"})
	//keys = append(keys, FormatType{IsInt: true, Content: "0"})
	//keys = append(keys, FormatType{Content: "a"})
	v := []interface{}{20, 30}
	m := UpdateJsonOnKeys(msg, v, keys)
	Log(m)
	k := GetJSON(m, keys)
	Log(k)
	value := map[string]interface{}{}
	GetAllJSON(msg, keys, value)
	Log(value)
	keys = BuildKeys("m.k[0].add")
	v = []interface{}{90}
	m = AddJSON(msg, v, keys)
	Log("增加字段add 90后的消息体：", m)
	keys = BuildKeys("m.k[5].add")
	v = []interface{}{999}
	m = AddJSON(msg, v, keys)
	Log("增加字段add 999后的消息体：", m)
	keys = BuildKeys("m.k[5].add")
	m = DeleteJSON(m, keys)
	Log("删除【5】字段add 999后的消息体：", m)
}

func TestUpdateJSON(t *testing.T) {
	s := `   { "m" : { "k": [ { "a" : 1, "b" : "c", "x": [1, 2, 3] }, { "a" : 2, "b" : "d", "x":[3,4,5] } ], "s": "hello" }, "w": "world"}`
	msg, fields := WalkJSON([]byte(s), "R")
	for k, v := range fields {
		Log("key", k, "value", v)
	}
	Log(msg)
	Log("+++++++++++++++++++++++++++++++++++")
	s2 := `   { "m" : { "k": [ { "a" : 10, "b" : "cx", "x": [10, 20] }, { "b" : "dx" }, { "c" : "ex", "x":[33,66] } ], "h": "hello world" }, "g": "hello world"}`
	msg2, fields2 := WalkJSON([]byte(s2), "R")
	Log(msg2)
	for k, v := range fields2 {
		Log("key", k, "value", v)
	}
	msg = UpdateJSON(msg, msg2)
	Log(msg)
}
