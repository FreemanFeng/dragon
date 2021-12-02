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

package builder

import (
	. "github.com/FreemanFeng/dragon/dragon/common"
)

func Run() {
	ch := AnyChannel(ModuleBuilder)
	for {
		x := <-ch
		k := x.(TaskRequest)
		go process(k)
		Debug(ModuleBuilder, k.Task, k.Project, k.Suite, k.Group, k.Case)
	}
}

func process(r TaskRequest) {
	ch := AnyChannel(Join(EMPTY, ModuleParser, r.Task))
	for {
		x := <-ch
		t := x.(*Testing)
		rewrite(r, t)
		go load(r, t)
		AnyChannel(Join(EMPTY, ModuleBuilder, r.Task)) <- t
		Debug(ModuleBuilder, t)
	}
}
