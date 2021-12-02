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

package templates

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/FreemanFeng/dragon/dragon/common"
	"github.com/FreemanFeng/dragon/dragon/util"
)

func Parse(ch chan TaskRequest) {
	for {
		r := <-ch
		Debug(ParserGlobal, "Request", r)
		start(r)
	}
}

func start(r TaskRequest) {
	root := ParserPath(TPath, r.Project, ParserTemplates)
	files := util.GetFiles(root, STAR)
	nodes := map[string]*TemplateType{}
	for _, k := range files {
		Debug(ParserTemplates, "root", root, "file", k)
		recognize(r, root, k, nodes)
	}
	AnyChannel(Join(EMPTY, ParserTemplates, r.Task)) <- nodes
}

func recognize(r TaskRequest, root, file string, nodes map[string]*TemplateType) {
	if IsDir(file) {
		return
	}
	name := filepath.Base(file)
	p := &TemplateType{Name: name, Files: map[string]string{}}
	h := strings.Split(name, DOT)
	n := len(h)
	for i, s := range h {
		if i == n-1 {
			continue
		}
		if i > 0 && IsAction(s) {
			p.Actions = append(p.Actions, s)
			continue
		}
		p.Keys = append(p.Keys, s)
	}
	nodes[file] = p
	t := strings.Split(file, root)
	s := strings.TrimPrefix(t[1], string(os.PathSeparator))
	for strings.Index(s, string(os.PathSeparator)) == 0 {
		s = strings.TrimPrefix(s, string(os.PathSeparator))
	}
	nodes[s] = p
	Debug(ModuleParser, "recognize", nodes[s].Files, "keys:", nodes[s].Keys, "actions:", nodes[s].Actions)
	t = strings.Split(s, string(os.PathSeparator))
	if len(t) < 2 {
		p.Files[DOT] = file
	} else {
		top := t[0]
		p.Files[top] = file
		s = t[1]
		k, ok := nodes[s]
		if !ok {
			nodes[s] = p
		} else {
			k.Files[top] = file
		}
		Debug(ModuleParser, "recognize", s, nodes[s].Files)
		k, ok = nodes[name]
		if !ok {
			nodes[name] = p
		} else {
			k.Files[top] = file
		}
		Debug(ModuleParser, "recognize", nodes[name].Files)
	}
}
