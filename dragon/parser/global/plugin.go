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
package global

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/FreemanFeng/dragon/dragon/util"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func ParsePlugin(r TaskRequest, start int, content [][]byte, nodes map[string]*PluginType) (int, int, string) {
	var calls []string
	var mutes []string
	var p *PluginType
	rw := regexp.MustCompile(`\w+`)
	rs := regexp.MustCompile(`\s+`)
	n := len(content)
	end := start
	id := EMPTY
	title := EMPTY
	desc := EMPTY
	mode := EMPTY
	name := EMPTY
	startup := EMPTY
	buildup := EMPTY
	proto := EMPTY
	port := EMPTY
	closed := 0
	extends := map[string]string{}
	for i := start; i < n; i++ {
		c := bytes.TrimSpace(content[i])
		// 遇到下一区域，跳出处理
		if MatchStart(c, MarkRegion) {
			break
		}
		if MatchStart(c, MarkRegion2) {
			s := string(bytes.TrimSpace(c))
			h := rs.Split(s, UNDEFINED)
			if len(h) < 2 {
				Err(ModuleParser, "获取不到插件标题")
				return end, FAILED, MsgIncompleteRule
			}
			title = h[1]
			Debug("插件标题为", title)
			continue
		}
		if len(title) > 0 && !MatchStart(c, MarkClass) {
			desc += string(c)
			b := bytes.TrimSpace(c)
			k := bytes.Split(b, []byte(OpEqual))
			if len(k) > 1 {
				s := string(bytes.TrimSpace(k[0]))
				// 处理扩展属性
				if strings.Index(s, SPACE) > 0 {
					x := strings.Split(s, SPACE)
					s = x[0]
				}
				t := string(bytes.TrimSpace(k[1]))
				h := rs.Split(t, UNDEFINED)
				switch s {
				case ModeCN, ModeEN:
					mode = h[0]
				case ProtoCN, ProtoEN:
					proto = h[0]
				case PortCN, PortEN:
					port = h[0]
				case ClosedCN, ClosedEN:
					x, err := strconv.Atoi(h[0])
					if err == nil && x > 0 {
						closed = 1
					}
				case NameCN, NameEN:
					name = h[0]
				case StartCN, StartEN:
					startup = string(k[1])
				case BuildCN, BuildEN:
					buildup = string(k[1])
				case CallsCN, CallsEN:
					calls = append(calls, h...)
				case MutesCN, MutesEN:
					mutes = append(mutes, h...)
				case ExtendCN, ExtendEN:
					x := bytes.Index(b, []byte(SPACE))
					n := len(b)
					k = bytes.Split(b[x:n], []byte(OpEqual))
					s = string(bytes.TrimSpace(k[0]))
					t = string(bytes.TrimSpace(k[1]))
					extends[s] = t
				}
			}
			continue
		}

		end = i
		if MatchStart(c, MarkClass) {
			id = string(rw.Find(c))
			Debug(ModuleParser, "Find Plugin ID", id)
			_, ok := nodes[id]
			if ok {
				Err("插件ID", id, "重复")
				return end, FAILED, MsgDuplicatedName
			}
			nodes[id] = &PluginType{PID: id, Title: title, Desc: desc, Mode: mode,
				Name: name, Start: startup, Build: buildup, Calls: calls, Mutes: mutes,
				Proto: proto, Port: port, Closed: closed,
				Funcs:   map[string]interface{}{},
				Bins:    map[string]string{},
				Extends: map[string][]string{}}
			p = nodes[id]
			for x, y := range extends {
				p.Extends[x] = rs.Split(y, UNDEFINED)
			}

			title = EMPTY
			desc = EMPTY
			mode = EMPTY
			name = EMPTY
			startup = EMPTY
			buildup = EMPTY
			proto = EMPTY
			closed = 0
			port = EMPTY
			calls = nil
			mutes = nil
			extends = map[string]string{}
			c = bytes.TrimSpace(bytes.TrimPrefix(c, []byte(MarkClass)))
			c = bytes.TrimSpace(bytes.TrimPrefix(c, []byte(id)))
			k := bytes.Split(c, []byte(VerticalBar))
			for _, b := range k {
				t := bytes.Split(b, []byte(OpEqual))
				if len(t) < 2 {
					continue
				}
				s := string(bytes.TrimSpace(t[1]))
				h := rs.Split(s, UNDEFINED)
				s = string(bytes.TrimSpace(t[0]))
				switch s {
				case CallsCN, CallsEN:
					p.Calls = append(p.Calls, h...)
				case MutesCN, MutesEN:
					p.Mutes = append(p.Mutes, h...)
				case ModeCN, ModeEN:
					p.Mode = h[0]
				case ProtoCN, ProtoEN:
					p.Proto = h[0]
				case PortCN, PortEN:
					p.Port = h[0]
				case ClosedCN, ClosedEN:
					x, err := strconv.Atoi(h[0])
					if err == nil && x > 0 {
						p.Closed = 1
					}
				case NameCN, NameEN:
					p.Name = h[0]
				}

			}
			build(r, p)
			Info(ModuleParser, "发现插件", id, "标题", p.Title, "模式", p.Mode, "名称", p.Name,
				"命令", p.Start, "构建", p.Build, "端口", p.Port,
				"调用", p.Calls, "屏蔽", p.Mutes, "协议", p.Proto, "扩展", p.Extends)
			continue
		}
		// 空行或注释
		if len(c) == 0 || MatchStart(c, MarkComment) || MatchStart(bytes.TrimSpace(c), MarkComment2) {
			continue
		}
		// 判断是否已获取插件ID
		_, ok := nodes[id]
		if id == EMPTY || !ok {
			Err("没有获取到插件ID")
			return end, FAILED, MsgIncompleteRule
		}
	}

	return end, SUCCESSFUL, MsgSuccessful
}

func build(r TaskRequest, p *PluginType) {
	up := p.Mode
	if p.Mode == RPC {
		up = SRC
	}
	p.Path = ParserPath(TPath, r.Project, ParserPlugins, up, p.Name)
	// 支持Python等编程语言
	if strings.Contains(p.Name, DOT) {
		p.Path = ParserPath(TPath, r.Project, ParserPlugins, up)
	}
	if p.Mode == SRC && !IsDir(p.Path) {
		return
	}

	// 插件已关闭
	if p.Closed == 1 {
		return
	}

	if p.Mode == BIN {
		for _, v := range p.Calls {
			p.Bins[v] = p.Path
		}
		Log("+++++++ 执行文件映射", p.Bins)
		return
	}

	// RPC服务需手动启
	if p.Mode == RPC {
		p.IsRPC = true
		initRPC(p)
		return
	}

	p.IsSRC = true

	// TODO windows DLL太诡异，放弃 - Freeman
	if runtime.GOOS == OSWindows || !IsGoCode(p) {
		p.IsRPC = true
		if p.Port == EMPTY {
			Err(ModuleParser, "!!!!!!插件端口未配置，配置格式为：端口=<端口号> 或 port=<端口号>")
			return
		}
		compileEXE(r, p)
		bin := getBin(p)
		if IsFile(bin) {
			loadEXE(r, p)
		} else {
			loadScript(r, p)
		}
		initRPC(p)
		return
	}
	compile(r, p)
	load(r, p)
}

func hasFile(p *PluginType, name string) bool {
	s := []string{p.Path, string(os.PathSeparator), name}
	h := strings.Join(s, EMPTY)
	return IsFile(h)
}

func getBin(p *PluginType) string {
	s := []string{p.Path, string(os.PathSeparator), p.Name}
	if runtime.GOOS == OSWindows {
		s = append(s, ExtEXE)
	}
	if runtime.GOOS != OSWindows && strings.Contains(p.Name, DOT) {
		return EMPTY
	}
	return strings.Join(s, EMPTY)
}

func getURLBase(port string) string {
	s := []string{ProtoHTTP, COLON, SLASH, SLASH, LOCALHOST, COLON, port}
	return strings.Join(s, EMPTY)
}

func loadEXE(r TaskRequest, p *PluginType) {
	s := []string{"cd ", p.Path, " && "}
	Info(ModuleParser, ">>>>> 加载", p.Name)
	if p.Start != EMPTY {
		s = append(s, p.Start)
	} else {
		s = append(s, p.Name)
		if runtime.GOOS == OSWindows {
			s = append(s, ExtEXE, SPACE)
		}
		s = append(s, " -p ", p.Port)
	}

	c := strings.Join(s, EMPTY)
	Info(ModuleParser, "执行命令ExE", c)
	if runtime.GOOS == OSWindows {
		p.Cmd = exec.Command("cmd", "/C", c)
	} else {
		p.Cmd = exec.Command("sh", "-c", c)
	}
	err := p.Cmd.Start()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		os.Exit(1)
	}
	//err = cmd.Wait()
	//if err != nil {
	//	Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
	//	os.Exit(1)
	//}
}

func loadScript(r TaskRequest, p *PluginType) {
	ext := path.Ext(p.Name)
	s := []string{"cd ", p.Path, " && "}
	Info(ModuleParser, ">>>>> 加载", p.Name)
	if p.Start != EMPTY {
		s = append(s, p.Start)
	} else {
		// 目前仅支持python/java
		switch ext {
		case ExtPy:
			s = append(s, "python3 ", p.Name, SPACE, p.Port)
		case ExtJar:
			s = append(s, "java -Dport=", p.Port, " -jar ", p.Name)
		default:
			return
		}
	}

	c := strings.Join(s, EMPTY)
	Info(ModuleParser, "执行脚本", c)
	if runtime.GOOS == OSWindows {
		p.Cmd = exec.Command("cmd", "/C", c)
	} else {
		p.Cmd = exec.Command("sh", "-c", c)
	}
	err := p.Cmd.Start()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		os.Exit(1)
	}
	//err = cmd.Wait()
	//if err != nil {
	//	Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
	//	os.Exit(1)
	//}
}

func initRPC(p *PluginType) {
	var msg []string
	s := []string{getURLBase(p.Port), PathInit}
	reqURL := strings.Join(s, EMPTY)
	b, code, _ := util.RequestHttp(reqURL, MethodGet, nil, nil, nil)
	for i := 1; code != http.StatusOK && i < MaxRetry; i++ {
		time.Sleep(2 * time.Second)
		b, code, _ = util.RequestHttp(reqURL, MethodGet, nil, nil, nil)
	}
	err := json.Unmarshal(b, &msg)
	if err != nil {
		Err(ModuleParser, "无法调用RPC服务Init接口", err)
		os.Exit(1)
	}
	for _, v := range msg {
		Info(ModuleParser, "初始化函数", v)
		// 若函数调用列表不为空，将过滤掉没有指定的函数
		if len(p.Calls) > 0 && !HasKey(v, p.Calls) {
			continue
		}
		// 若函数屏蔽列表不为空，将过滤掉命中的函数
		if len(p.Mutes) > 0 && HasKey(v, p.Mutes) {
			continue
		}
		p.Funcs[v] = RunRPC
	}
}

func RunRPC(port, service string, buf []byte) []byte {
	s := []string{getURLBase(port), PathRun, SLASH, service}
	reqURL := strings.Join(s, EMPTY)
	b, code, _ := util.RequestHttp(reqURL, MethodPost, nil, nil, buf)
	if code != http.StatusOK {
		Err(ModuleParser, ">>>>> 调用RPC服务", service, "出错，状态码", code)
		os.Exit(1)
	}
	return b
}

func killProcess(r TaskRequest, p *PluginType) {
	var cmd *exec.Cmd
	s := []string{p.Name}
	bin := getBin(p)
	if !IsFile(bin) {
		return
	}
	if runtime.GOOS == OSWindows {
		s = append(s, ExtEXE)
	}
	c := strings.Join(s, EMPTY)
	if runtime.GOOS == OSWindows {
		cmd = exec.Command("taskkill", "/F", "/T", "/IM", c)
	} else {
		s = []string{"pgrep", c, "|", "xargs", "kill", "-9"}
		c = strings.Join(s, SPACE)
		cmd = exec.Command("sh", "-c", c)
	}
	err := cmd.Run()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		os.Exit(1)
	}
}

func killPort(r TaskRequest, p *PluginType) {
	var cmd *exec.Cmd
	if p.Port == EMPTY {
		return
	}
	if runtime.GOOS == OSWindows {
		s := []string{"netstat -ano | find ", "\"", p.Port, "\""}
		c := strings.Join(s, EMPTY)
		cmd = exec.Command("cmd", "/C", c)
	} else {
		s := []string{"lsof -i:", p.Port}
		c := strings.Join(s, EMPTY)
		cmd = exec.Command("sh", "-c", c)
	}
	out, err := cmd.Output()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令出错:", err)
		os.Exit(1)
	}
	Info("获取进程ID", out)
}

func checkNewer(r TaskRequest, p *PluginType, bin string) bool {
	fs := util.GetFiles(p.Path, STAR)
	for _, f := range fs {
		if IsNewer(f, bin) {
			return true
		}
	}
	return false
}

func compileEXE(r TaskRequest, p *PluginType) {
	var cmd *exec.Cmd
	bin := getBin(p)
	if IsFile(bin) && !checkNewer(r, p, bin) || bin == EMPTY {
		return
	}
	killProcess(r, p)
	Info(ModuleParser, ">>>> 重新编译", bin)
	s := []string{"cd ", p.Path, " && "}
	if p.Build != EMPTY {
		s = append(s, p.Build)
	} else if IsGoCode(p) {
		s = append(s, "go build")
	} else if hasFile(p, MF) || hasFile(p, MF2) || hasFile(p, MF3) {
		s = append(s, "make")
	} else {
		return
	}
	c := strings.Join(s, EMPTY)
	Info(ModuleParser, "执行命令", c)
	if runtime.GOOS == OSWindows {
		cmd = exec.Command("cmd", "/C", c)
	} else {
		cmd = exec.Command("sh", "-c", c)
	}
	err := cmd.Run()
	if err != nil {
		Err(ModuleParser, "!!!!! 执行命令", c, "出错:", err)
		return
	}
	Info(ModuleParser, "执行命令", c, "成功!")
}
