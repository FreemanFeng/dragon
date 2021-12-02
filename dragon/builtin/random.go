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

package builtin

import (
	"math/rand"
	"strconv"

	. "github.com/nuxim/dragon/dragon/common"

	"github.com/google/uuid"
)

func BuildRandom(params ...string) string {
	n := len(params)
	selection := CHARS
	cnt := DefaultCharsLength
	switch {
	case n == 1:
		selection = params[0]
	case n > 1:
		selection = params[0]
		cnt, _ = strconv.Atoi(params[1])
	}
	x := []byte(selection)
	n = len(x)
	b := make([]byte, cnt)
	for i := 0; i < cnt; i++ {
		k := rand.Intn(n)
		b[i] = x[k]
	}
	return string(b)
}

func RandomString(params ...string) string {
	n := len(params)
	cnt := DefaultCharsLength
	if n > 0 {
		cnt, _ = strconv.Atoi(params[0])
	}
	return BuildRandom(CHARS, strconv.Itoa(cnt))
}

func Random(params ...string) string {
	return RandomString(params...)
}

func RandomDigit(params ...string) string {
	n := len(params)
	cnt := DefaultCharsLength
	if n > 0 {
		cnt, _ = strconv.Atoi(params[0])
	}
	return BuildRandom(DIGITS, strconv.Itoa(cnt))
}

func UUID(params ...string) string {
	k, err := uuid.NewUUID()
	if err != nil {
		return EMPTY
	}
	return k.String()
}
