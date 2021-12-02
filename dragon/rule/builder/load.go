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
	"bytes"
	"io/ioutil"
	"path"
	"strings"

	. "github.com/nuxim/dragon/dragon/common"
)

func load(r TaskRequest, t *Testing) {
	hosts := selectHosts(r, t)
	id := UNDEFINED
	n := len(hosts)
	max := n
	if r.IsStateless {
		max = 1
	}
	// 首先加载前置条件后置条件，以便执行用例前可以跑
	loadTesting(r.Task, t)
	for k := range t.Suites {
		id = loadSuite(r.Task, k, hosts, n, max, id, t, t.Suites[k])
	}
}

func selectHosts(r TaskRequest, t *Testing) []string {
	// 过滤掉相同的host
	m := map[string]int{}
	var hosts []string
	for _, v := range r.Hosts {
		m[v] = 1
	}
	if !r.IsExtraHosts && len(m) > 0 {
		for k := range m {
			hosts = append(hosts, k)
		}
		return hosts
	}
	c := t.Config
	s, ok := c.Settings[StaticHosts]
	if ok {
		for _, v := range s {
			m[v] = 1
		}
	}
	for k := range m {
		hosts = append(hosts, k)
	}
	if len(hosts) == 0 {
		return []string{LOCALHOST}
	}
	return hosts
}

func loadTesting(task string, t *Testing) {
	gd := &GroupData{Scope: ScopeAll, Cases: map[string][]*VarData{}, CFs: map[string]*VarData{}, FDs: map[string]*VarData{}}
	m := map[string]*FileType{}
	for _, c := range t.Setup {
		gd.CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
			Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
			Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
			VPos: map[string]int{}, Ctx: map[string]interface{}{}}
		d := gd.CFs[c.CID]
		loadControlFlowData(task, t, c, m, d)
	}
	for _, c := range t.Teardown {
		gd.CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
			Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
			Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
			VPos: map[string]int{}, Ctx: map[string]interface{}{}}
		d := gd.CFs[c.CID]
		loadControlFlowData(task, t, c, m, d)
	}
	for _, c := range t.ControlFlow {
		gd.CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
			Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
			Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
			VPos: map[string]int{}, Ctx: map[string]interface{}{}}
		d := gd.CFs[c.CID]
		loadControlFlowData(task, t, c, m, d)
	}
	for _, c := range t.Checks {
		gd.CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
			Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
			Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
			VPos: map[string]int{}, Ctx: map[string]interface{}{}}
		d := gd.CFs[c.CID]
		loadControlFlowData(task, t, c, m, d)
	}
	AnyChannel(Join(EMPTY, ModuleRunner, task)) <- gd
}

func loadSuite(task, suite string, hosts []string, n, max, id int, t *Testing, p *TestSuite) int {
	gs := make([]*GroupData, max)
	for i := 0; i < max; i++ {
		gs[i] = &GroupData{SID: suite, Cases: map[string][]*VarData{},
			CFs: map[string]*VarData{}, FDs: map[string]*VarData{}}
		m := map[string]*FileType{}
		for _, c := range t.Setup {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range t.Teardown {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range t.ControlFlow {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range t.Checks {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range t.Config.Flows {
			gs[i].FDs[c.FID] = &VarData{ID: c.FID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].FDs[c.FID]
			loadFlowData(task, t, c, m, d)
		}
		for _, c := range p.Setup {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range p.Teardown {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range p.ControlFlow {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range p.Checks {
			gs[i].CFs[c.CID] = &VarData{ID: c.CID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].CFs[c.CID]
			loadControlFlowData(task, t, c, m, d)
		}
		for _, c := range p.Flows {
			gs[i].CFs[c.FID] = &VarData{ID: c.FID, Title: c.Title, Data: map[string][]*FileType{}, Bins: map[string]string{},
				Tops: map[string][]string{}, Msgs: map[string]interface{}{}, Funcs: map[string]*PluginType{},
				Tags: map[string]*OpType{}, Keys: map[string][]FormatType{}, Vars: map[string]interface{}{},
				VPos: map[string]int{}, Ctx: map[string]interface{}{}}
			d := gs[i].FDs[c.FID]
			loadFlowData(task, t, c, m, d)
		}
	}
	for k, v := range p.Groups {
		id = (id + 1) % n
		gid := id % max
		gs[gid].GID = k
		gs[gid].Host = hosts[id]
		m := map[string]*FileType{}
		// 加载用例数据
		for _, c := range v.Cases {
			for i := 0; i < c.Concurrence; i++ {
				cd := &VarData{ID: c.TID, Title: c.Title, Host: hosts[id],
					Data: map[string][]*FileType{}, Tags: map[string]*OpType{}, Funcs: map[string]*PluginType{},
					Tops: map[string][]string{}, Msgs: map[string]interface{}{},
					Keys: map[string][]FormatType{}, Bins: map[string]string{}, Vars: map[string]interface{}{},
					VPos: map[string]int{}, Ctx: map[string]interface{}{}}
				loadGroupData(task, t, v, cd)
				loadCaseData(task, t, c, m, cd)
				gs[gid].Cases[c.TID] = append(gs[gid].Cases[c.TID], cd)
			}
		}

		AnyChannel(Join(EMPTY, ModuleRunner, task)) <- gs[gid]
	}
	return id
}
func loadGroupData(task string, t *Testing, p *GroupType, cd *VarData) {
	loadFuncs(p.Plugins, t, cd)
	loadBins(p.Plugins, t, cd)
}

func loadCaseData(task string, t *Testing, p *CaseType, m map[string]*FileType, cd *VarData) {
	loadFuncs(p.Plugins, t, cd)
	loadBins(p.Plugins, t, cd)
	for k, v := range p.Vars {
		loadTemplates(task, k, t.Templates, v, m, cd.Data)
	}
	for k := range cd.Data {
		fixTemplates(task, k, t.Templates, m, cd.Data)
	}
	for k, v := range p.Vars {
		// 加载上下文
		loadContexts(cd, k, v)
	}
	for _, v := range p.Steps {
		sp := &v.Setting
		if sp.IsTag {
			loadTags(cd, p.Vars, sp)
		} else {
			loadOpsKey(cd, p.Vars, sp)
		}
		// 加载校验keys
		loadExpectKeys(cd, p.Vars, v.Expects)
	}
	for k := range m {
		recognize(t, m[k], cd)
	}
	if p.Rounds > 0 && cd.Total > int(p.Rounds) {
		cd.Total = int(p.Rounds)
	}
}

func loadFlowData(task string, t *Testing, p *FlowType, m map[string]*FileType, cd *VarData) {
	loadFuncs(p.Plugins, t, cd)
	loadBins(p.Plugins, t, cd)
	for k, v := range p.Vars {
		loadTemplates(task, k, t.Templates, v, m, cd.Data)
	}
	for k := range cd.Data {
		fixTemplates(task, k, t.Templates, m, cd.Data)
	}
	for k, v := range p.Vars {
		// 加载上下文
		loadContexts(cd, k, v)
	}
	for _, v := range p.Steps {
		sp := &v.Setting
		if sp.IsTag {
			loadTags(cd, p.Vars, sp)
		} else {
			loadOpsKey(cd, p.Vars, sp)
		}
		// 加载校验keys
		loadExpectKeys(cd, p.Vars, v.Expects)
	}
	for k := range m {
		recognize(t, m[k], cd)
	}
}

func loadControlFlowData(task string, t *Testing, p *ControlType, m map[string]*FileType, cd *VarData) {
	loadFuncs(p.Plugins, t, cd)
	loadBins(p.Plugins, t, cd)
	for k, v := range p.Vars {
		loadTemplates(task, k, t.Templates, v, m, cd.Data)
	}
	for _, v := range p.Flows {
		sp := &v.Setting
		if sp.IsTag {
			loadTags(cd, p.Vars, sp)
		} else {
			loadOpsKey(cd, p.Vars, sp)
		}
	}
	for k := range m {
		recognize(t, m[k], cd)
	}
}

func loadTemplates(task, key string, t map[string]*TemplateType, vp *VarType, m map[string]*FileType, g map[string][]*FileType) {
	for _, v := range vp.Templates {
		//Log(ModuleBuilder, ">>>>>> 模板名称", v.Key, "key", key)
		k, ok := t[v.Key]
		if !ok {
			//Log(ModuleBuilder, ">>>>>> 模板", v.Key, "不存在, key", key)
			_, ok = g[v.Key]
			if ok {
				//Log(ModuleBuilder, ">>>>>> 已存在", v.Key, "对应的模板数据,key", key)
				g[key] = append(g[key], &FileType{Name: v.Key, IsVar: true})
			}
			continue
		}
		for top, kp := range k.Files {
			_, ok = m[kp]
			if !ok {
				c, e := ioutil.ReadFile(kp)
				if e != nil {
					Debug(ModuleBuilder, e)
					continue
				}
				b := checkBin(c, kp)
				//Log(ModuleBuilder, ">>>>>> 文件路径", path, "Name", v.Key, "Top", top, "ID", key)
				m[kp] = &FileType{Content: c, Top: top, Name: v.Key, Path: kp, IsBin: b}
			}
			g[key] = append(g[key], m[kp])
		}
		//Log(ModuleBuilder, ">>>>>> 模板key", key, "包含模板文件:")
		//for _, v := range g[key] {
		//	Log(ModuleBuilder, ">>>>>>+++++ 模板路径", v.Path, "模板名称", v.Name, "top", v.Top)
		//}
	}
}

func fixTemplates(task, key string, t map[string]*TemplateType, m map[string]*FileType, g map[string][]*FileType) {
	if len(g[key]) == 0 || g[key][0].IsVar {
		return
	}
	p := g[key][0]
	k, ok := t[p.Name]
	s := strings.Split(p.Name, DOT)
	n := len(s) - 1
	// 最大匹配
	//for k := range t {
	//	Log("++++++++++ 模板文件名：", k)
	//}
	for n > 1 {
		n--
		h := strings.Join(s[:n], DOT)
		x := []string{h, ExtJson}
		name := strings.Join(x, EMPTY)
		Log(ModuleBuilder, ">>>>>>>>>> 匹配文件名", name, "n:", n, "s[n-1]", s[n-1])
		k, ok = t[name]
		if !ok && IsAction(s[n-1]) {
			Log(ModuleBuilder, ">>>>>>>>>> 退出不匹配", p.Name, "n:", n, "s[n-1]", s[n-1])
			break
		} else if !ok {
			continue
		}
		tops := map[string]int{}
		for _, v := range g[key] {
			tops[v.Top] = 1
		}
		for top, kp := range k.Files {
			_, ok = tops[top]
			if ok {
				continue
			}
			g[key] = append(g[key], m[kp])
			Log(ModuleBuilder, ">>>>>>> 成功合并文件", kp)
		}
	}
}

func loadContexts(cd *VarData, key string, vt *VarType) {
	for _, v := range vt.Templates {
		if !v.IsCtx {
			continue
		}
		for _, v := range cd.Data[v.ID] {
			msg := ToJSON(v.Content)
			cd.Ctx[key] = UpdateJSON(cd.Ctx[key], msg)
		}
	}
}

// 部分识别Bin
func checkBin(c []byte, p string) bool {
	ext := path.Ext(p)
	// 假设以下后缀名的文件都不是可执行文件
	//Log(ModuleBuilder, "识别文件", p, "后缀名", ext)
	if Equal(ext, ExtJson, ExtTxt, ExtLog, ExtCfg, ExtIni,
		ExtGo, ExtPy, ExtJava, ExtMd, ExtBash, ExtSh) {
		return false
	}
	n := len(c)
	if n > MaxBinRead {
		n = MaxBinRead
	}
	for i := 0; i < n; i++ {
		if int(c[i]) > 127 {
			return true
		}
	}
	return false
}

// 映射插件函数到运行时数据
func loadFuncs(plugins []string, t *Testing, cd *VarData) {
	for k, v := range t.Config.Plugins {
		if HasKey(k, plugins) {
			for name, _ := range v.Funcs {
				cd.Funcs[name] = v
			}
		}
	}
}

// 映射插件函数到运行时数据
func loadBins(plugins []string, t *Testing, cd *VarData) {
	for k, v := range t.Config.Plugins {
		if HasKey(k, plugins) {
			for name, path := range v.Bins {
				cd.Bins[name] = path
			}
		}
	}
}

// Tag对应的值有可能是消息字段或其他变量，需要映射
func loadTags(cd *VarData, vt map[string]*VarType, sp *SettingType) {
	for i := range sp.Ops {
		p := &sp.Ops[i]
		if p.Op != OpEqual || p.Key == EMPTY {
			continue
		}
		cd.Tags[p.Key] = p
		for _, s := range p.Value {
			// 如果是字段赋值，需要加入字段映射
			LoadKeys(cd, vt, s)
		}
	}
}

// 仅仅是映射相应消息字段
func loadOpsKey(cd *VarData, vt map[string]*VarType, st *SettingType) {
	n := len(st.Ops)
	if n > cd.Total {
		cd.Total = n
	}
	for _, v := range st.Ops {
		if st.IsInt || st.IsList {
			continue
		}
		LoadKeys(cd, vt, v.Key)
		for _, s := range v.Value {
			// 如果是字段赋值，需要加入字段映射
			LoadKeys(cd, vt, s)
		}
		// 上下文赋值
		s := strings.Split(v.Key, DOT)
		key := s[0]
		k, ok := vt[key]
		if !ok {
			continue
		}
		// 预设取值范围，便于后续判断是否为上下文并从取值范围中取具体值
		if k.IsFlow {
			cd.Ctx[v.Key] = v.Value
		}
	}
}

func loadExpectKeys(cd *VarData, vt map[string]*VarType, ops []OpType) {
	n := len(ops)
	if n > cd.Total {
		cd.Total = n
	}
	for _, v := range ops {
		LoadKeys(cd, vt, v.Key)
		for _, s := range v.Value {
			// 如果是字段赋值，需要加入字段映射
			LoadKeys(cd, vt, s)
		}
	}
}

// 识别JSON格式或带动态标签的消息模板
func recognize(t *Testing, p *FileType, cd *VarData) {
	if p.IsBin {
		return
	}
	b := bytes.TrimSpace(p.Content)
	// 找到包含动态标签的消息模板
	for k := range cd.Tags {
		if MatchInside(b, k) {
			//Log(ModuleCommon, p.Name, "包含tag", k)
			p.Tags = append(p.Tags, k)
		}
	}
	// 部分消息模板使用动态标签，所以是无法通过JSON序列化的，这部分模板无法识别，待后续识别
	if MatchStart(b, LeftBracket) && MatchEnd(b, RightBracket) {
		p.IsList = true
		if IsJsonList(b) {
			p.IsJSON = true
			Debug(ModuleBuilder, p.Name, "是JSON列表")
		}
	}
	if MatchStart(b, LeftBrace) && MatchEnd(b, RightBrace) {
		p.IsMap = true
		if IsJsonMap(b) {
			p.IsJSON = true
			Debug(ModuleBuilder, p.Name, "是JSON字典")
		}
	}
}
