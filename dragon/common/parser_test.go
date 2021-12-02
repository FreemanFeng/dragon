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

func TestParseSetting(t *testing.T) {
	var c []byte
	c = []byte("k -> name = a b c id = 1 2 3 x = 3 4 5")
	ParseSetting(c)
	c = []byte("A.x = (1+2)*5 3>>2 1<<3 5|6 5~6 5&6 5**6 5//6 7/8 7*8 5^6 3..90,B.x..D.k,1,2,3 C...a 40..50 7,8,9 {} [] on a = 1 2 3 b = 4 5 6")
	ParseSetting(c)
	c = []byte("D -- c d e on a = 1 2 3 b = 4 5 6 c = 7 8 9")
	ParseSetting(c)
	c = []byte("3.9:A...t = D...s with x y z on A...x = 5 6 7 D...y = 6 8 9")
	ParseSetting(c)
	c = []byte("3.9:A...t = D...s on A...x = 5 6 7 D...y = 6 8 9 with x y z ")
	ParseSetting(c)
	c = []byte("3.9:A...t = 6 8 9 with x y z on A...x = 5 6 7 D...y = D...s")
	ParseSetting(c)
	c = []byte("3..9:A...t = 6 8 9 B...x = 3 4 5 on A...x = 5 6 7 D...y = D...s with x y z")
	ParseSetting(c)
	c = []byte("a = 2  b = 3 on c = 4 5 6 d = 7 8 9")
	ParseSetting(c)
	c = []byte("a = 2 3 4  b = 3 4 5 on c = 4 5 6 d = 7 8 9")
	ParseSetting(c)
	c = []byte("C ++ a = 2 3 4  b = 3 4 5 on c = 4 5 6 d = 7 8 9")
	ParseSetting(c)
	c = []byte("C ++ a = 2  b = 3 on c = 4 5 6 d = 7 8 9")
	ParseSetting(c)
	c = []byte("5..10 : E.k - D.j = 300 or A.x - B.y = 20 on C.k = 300")
	ParseSetting(c)
	c = []byte("g -> A : a.json B : 10*e*.json:1..9.33 3 C : f*g*.json")
	ParseSetting(c)
	c = []byte("x = MD5 C...c")
	ParseSetting(c)
	c = []byte("A...message with 暴雨 on city = 广州 district = 白云")
	ParseSetting(c)
	c = []byte("A.code = 400 on A.city = \"\" or A.interval = 分钟")
	ParseSetting(c)
	c = []byte("A -- city")
	ParseSetting(c)
	c = []byte("A ++ hello=world")
	ParseSetting(c)
}

func TestBuildKeys(t *testing.T) {
	s := "A...message"
	k := BuildKeys(s)
	Log("s:", s, "k:", k)
	s = "A...a[0]...id"
	k = BuildKeys(s)
	Log("s:", s, "k:", k)
}
