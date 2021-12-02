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
	"fmt"
	"reflect"
	"testing"
)

func TestMatchEnd(t *testing.T) {
	c := []byte("[abc]")
	if MatchEnd(c, RightBracket) {
		print("matched []")
	}
	c = []byte("abc")
	if MatchEnd(c, RightBracket) {
		print("matched 2 []")
	}
	c = []byte("{abc}")
	if MatchEnd(c, RightBrace) {
		print("matched  {}")
	}
}

func TestStrMatchEnd(t *testing.T) {
	c := "__abc__"
	if StrMatchEnd(c, DoubleUnderScope) {
		print("match __")
	}
	if StrMatchStart(c, DoubleUnderScope) {
		print("match __ start")
	}
}

func TestMapValue(t *testing.T) {
	m := map[string]int{"k": 1, "a": 2}
	k, _ := m["x"]
	fmt.Println("不存在int时默认值", k)
	s := map[string]string{"k": "a", "a": "2"}
	x, _ := s["x"]
	fmt.Println("不存在str时默认值", x)
}

func Foo(m, c, r map[string]string) {
	fmt.Println(m, c, r)
	m["hello"] = "world!"
}

func TestFoo(t *testing.T) {
	m := map[string]interface{}{
		"test": Foo,
	}
	k := reflect.ValueOf(m["test"])
	in := make([]reflect.Value, 3)
	params := make([]map[string]string, 3)
	params[0] = map[string]string{
		"a": "1",
	}
	params[1] = map[string]string{
		"b": "2",
	}
	params[2] = map[string]string{
		"c": "3",
	}
	for i, v := range params {
		in[i] = reflect.ValueOf(v)
		fmt.Println(in[i])
	}
	k.Call(in)
	for i := range in {
		fmt.Println(in[i])
	}
}
