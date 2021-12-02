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
	"encoding/json"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinKeys(t *testing.T) {
	r := assert.New(t)
	r.Equal("a:b", Join(COLON, "a", "b"), EMPTY)
	r.Equal("a:2", JoinKeys(COLON, "a", 2), EMPTY)
	r.Equal("a:true", JoinKeys(COLON, "a", true), EMPTY)
	r.Equal("a:b", JoinKeys(COLON, "a", []byte("b")), EMPTY)
	r.Equal("a:1.000000:true:b", JoinKeys(COLON, "a", 1.0, true, []byte("b")), EMPTY)
}

func TestToString(t *testing.T) {
	r := assert.New(t)
	r.Equal("[1 2 3]", ToString([]int{1, 2, 3}), EMPTY)
}

func TestSortMap(t *testing.T) {
	r := assert.New(t)
	m := map[string]int{}
	m["kk"] = 1
	m["tt"] = 2
	m["aa"] = 3
	s := SortMap(m)
	Log(s)
	r.Equal("aa:3 kk:1 tt:2", s, EMPTY)
}

func TestNormalizeRegexp(t *testing.T) {
	text := "\\ . + * ? ( ) | [ ] { } ^ $"
	s := NormalizeRegexp(text)
	Log(s)
}

func TestReplaceSpaceAroundEqual(t *testing.T) {
	c := TrimSpaces([]byte(`a = 1  2  3 c = "4   5 6" 7 8 with x y z has a b c on k = 3 4 5 s = 7 8 9`))
	Log(string(c))
}

func Split(c []byte) {
	Log(string(c))
	var ot []OpType
	k := bytes.Split(c, []byte(SPACE))
	for _, v := range k {
		ot = SplitOps(v, ot)
	}
	for i := range ot {
		v := ot[i].Value
		n := len(v)
		if n == 1 && len(v[0]) == 2 && StrMatchStart(v[0], LeftBrace) && StrMatchEnd(v[0], RightBrace) {
			ot[i].Value = nil
			ot[i].IsMap = true
			Log("替换后", ot[i].Value, "是否字典", ot[i].IsMap)
		}
		if n < 2 {
			continue
		}
		if v[0] == LeftBracket && v[n-1] == RightBracket {
			ot[i].Value = v[1 : n-1]
			ot[i].IsList = true
			Log("替换后", ot[i].Value, "是否列表", ot[i].IsList)
		}
	}
	Log(ot)
}

func TestSplitOps(t *testing.T) {
	//c := TrimSpaces([]byte(" a  >  b   c < d   e has f 3 4  k without o 5 6  p = q  r != s x with 7 8 on k = 3 s = 5"))
	//Split(c)
	//c = TrimSpaces([]byte(" A...x = C...d on k = 3 s = 5 t = 8"))
	//Split(c)
	//c = TrimSpaces([]byte(" A...x = C...d with x y z on k = 3 s = 5 t = 10"))
	//Split(c)
	//c := TrimSpaces([]byte("A...t - D...s = 3 on k = 3 A...x = 5 6 7 D...y = 6 8 9"))
	//Split(c)
	//c = TrimSpaces([]byte("A...t = D...s with x y z on A...x = 5 6 7 D...y = 6 8 9 F"))
	//Split(c)
	c := TrimSpaces([]byte("A.x = (1+2)*5 3>>2 1<<3 5|6 5~6 5&6 5**6 5//6 7/8 7*8 5^6 3..90,B.x..D.k,1,2,3 C...a 40..50 7,8,9 {} [] on a = 1 2 3 b = 4 5 6"))
	Split(c)
	//c = TrimSpaces([]byte("D -- c d e on a = 1 2 3 b = 4 5 6"))
	//Split(c)
	//c = TrimSpaces([]byte("3..9.20 C ++ a = 2 4 5  b = 3 6 7 on c = 4 5 6 d = 7 8 9"))
	//Split(c)
	//c = TrimSpaces([]byte("A.d[0].x.a[1].b = 5 6 8..20 B.d[0].x.a[1].b = 5..25 7 9"))
	//Split(c)
	//c = TrimSpaces([]byte("A.c[2]...d[3]...k = 6..9 20 B.c[2]...d[3]...k = 6..9 8 20"))
	//Split(c)
	//c = TrimSpaces([]byte("C.x = MD5 A.a B.b C.c on A.x = 5..9 20 B.k = 2..20"))
	//Split(c)
	//c = TrimSpaces([]byte("MD5 a b c"))
	//Split(c)
	//c = TrimSpaces([]byte("x = {}"))
	//Split(c)
	//c = TrimSpaces([]byte("k -> 3..9 6..20"))
	//Split(c)
	//c = TrimSpaces([]byte("y = [4 5 6 9..100]"))
	//Split(c)
	//c = TrimSpaces([]byte("k -> name = a b c id = 1 2 3"))
	//Split(c)
	//c = TrimSpaces([]byte("g -> A : a.json B : 10*e*.json:1..9.33 3 C : f*g*.json"))
	//Split(c)
	//c = TrimSpaces([]byte("3.9 : V1 A.a B.b C.c D.d D...x"))
	//Split(c)
	c = TrimSpaces([]byte("3..9.30.99 : B.a = 3"))
	Split(c)
	//c = TrimSpaces([]byte("A...t = D...s with x y z"))
	//Split(c)
	//c = TrimSpaces([]byte("3.9:A...t = D...s with x y z on A...x = 5 6 7 D...y = 6 8 9"))
	//Split(c)
	c = TrimSpaces([]byte("F = GetWeatherFlow <- A:GetWeather.json:1 R:GetWeatherRule:2 K:GetWeatherRule2:4"))
	Split(c)
	c = TrimSpaces([]byte("1.2.3: F = GetWeatherFlow <- GetWeatherRule:2 GetWeatherRule2:4"))
	Split(c)
	c = TrimSpaces([]byte("1.2: G : GetWeatherFlow <- GetWeatherRule:2 GetWeatherRule2:4"))
	Split(c)
	c = TrimSpaces([]byte("next <- B:GetWeather C:GetWeather3"))
	Split(c)
	c = TrimSpaces([]byte("C.a.nonce = ( A.nonce + 1 ) * 2 B.nonce"))
	Split(c)
	c = TrimSpaces([]byte("G : GetWeatherFlow <- C:GetWeather.json:2 X:GetWeatherRule GetWeatherRule2:4"))
	Split(c)
}

func TestConvertSpacesInsideDoubleQuotes(t *testing.T) {
	c := []byte("'kkk   -- sss\" ++ \"ttt ++   aaa'")
	r := ConvertSpacesInsideQuotes(c)
	Log(string(c))
	Log(string(r))
}

func TestConvertByte(t *testing.T) {
	s := ""
	h := []string{s, s, s, s}
	k := strings.Join(h, "")
	b := []byte(k)
	Log("空字符值", b)
}

func TestJsonList(t *testing.T) {
	var m []interface{}
	var v [][]byte
	m = append(m, []byte("1"), []byte("hello"))
	b := JsonListByte(m...)
	json.Unmarshal(b, &v)
	Log("转换前m：", m)
	Log("转换后v：", v)
}

func TestGetFileName(t *testing.T) {
	k := "/a/b/c.go"
	Log(path.Base(k))
}

func TestSplitVars(t *testing.T) {
	s := "(A.nonce+B.a.nonce+1)*2"
	rs := SplitVars(s)
	Log(rs)
	Log(IsOpVars(s))
}
