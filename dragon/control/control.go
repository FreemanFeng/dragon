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

package control

import (
	"bytes"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nuxim/dragon/dragon/runner"

	"github.com/nuxim/dragon/dragon/proto"
	"github.com/nuxim/dragon/dragon/rule/builder"
	"github.com/nuxim/dragon/dragon/store"

	"github.com/nuxim/dragon/dragon/builtin"

	"github.com/nuxim/dragon/dragon/util"

	"github.com/gorilla/websocket"

	"github.com/nuxim/dragon/dragon/report"

	"github.com/nuxim/dragon/dragon/control/plugin"

	"github.com/nuxim/dragon/dragon/parser"

	"github.com/gin-gonic/gin"
	. "github.com/nuxim/dragon/dragon/common"
)

func Run(port int, debug bool) {
	rand.Seed(time.Now().Unix())
	InitLogConfig()
	go store.Run()
	InitGlobalConfig()
	if debug {
		go Profiling()
		IntChannel(FuncProfiling) <- GetFreePort()
	}
	go builtin.Run()
	go proto.Run()
	go parser.Run()
	go report.Run()
	go plugin.Run()
	go builder.Run()
	go runner.Run()
	Debug(ModuleControl, "is running on port", port)
	// 初始化引擎
	engine := gin.Default()
	// 注册路由和处理函数
	engine.GET(PathTestProject, TestProject)
	engine.GET(PathStopProject, StopProject)
	engine.GET(PathTestReport, RunTestReport)
	engine.GET(PathCallPlugin, CallPlugin)

	engine.GET(PathWSTask, RunWSTask)
	engine.GET(PathConfig, Config)
	engine.GET(PathStopShow, StopTask)
	// 绑定端口，然后启动应用
	s := []string{COLON, strconv.Itoa(port)}
	engine.Run(strings.Join(s, EMPTY))
}

func TestProject(context *gin.Context) {
	project := context.Param(ParamProject)
	ts := context.Query(ParamSuite)
	// 用,表示或的关系，用+表示且的关系，如 a+b,c表示名称含a且含b，或者含c
	group := context.Query(ParamGroup)
	tc := context.Query(ParamCase)
	// 用,表示或的关系，用+表示且的关系，如 a+b,c表示标题/ID含a且含b，或者含c
	tags := context.Query(ParamTags)
	p := context.Query(ParamConcurrence)
	concurrence, _ := strconv.Atoi(p)
	pd := context.Query(ParamDuration)
	duration := util.CalculateSecs(pd)
	hosts := context.Query(ParamHosts)
	isExtra := false
	if hosts == EMPTY {
		hosts = context.Query(ParamExtraHosts)
		isExtra = true
		Debug(ModuleControl, "Extra Hosts", hosts)
	}
	isStateless := false
	if concurrence > 0 {
		isStateless = true
	}
	smoking := context.Query(ParamSmoking)
	random := context.Query(ParamRandom)
	nocheck := context.Query(ParamNoCheck)

	task := MD5String(TPath, context.ClientIP(), project, p, hosts, pd,
		smoking, random, nocheck, util.GetMillisecond())
	r := TaskRequest{IsMocking: false, IsStateless: isStateless, IsExtraHosts: isExtra, Task: task, Project: project,
		Concurrence: concurrence, Duration: duration, Quit: make(chan int)}
	if hosts != EMPTY {
		r.Hosts = strings.Split(hosts, COMMA)
	}
	if ts != EMPTY {
		r.Suite = strings.Split(ts, COMMA)
	}
	if tc != EMPTY {
		r.Case = strings.Split(tc, COMMA)
	}
	if tags != EMPTY {
		r.Tags = strings.Split(tags, COMMA)
	}
	if group != EMPTY {
		r.Group = strings.Split(group, COMMA)
	}
	if smoking != EMPTY {
		r.IsSmoking = true
	}
	if random != EMPTY {
		r.IsRandom = true
	}
	if nocheck != EMPTY {
		r.IsNoCheck = true
	}
	start(context, r)
}

func StopProject(context *gin.Context) {
	path := context.Param(ParamPath)
	project := context.Param(ParamProject)
	task := context.Query(ParamTask)
	r := TaskRequest{Task: task, Project: project}
	stop(context, r, path)
}

func start(context *gin.Context, r TaskRequest) {
	AnyChannel(ModuleParser) <- r
	AnyChannel(ModuleBuilder) <- r
	AnyChannel(ModuleRunner) <- r
	AnyChannel(ModuleReport) <- r
	SetProjectTask(TPath, r.Project, r.Task, r.Quit)
	context.Redirect(http.StatusFound, strings.Join([]string{ParamReport, r.Task}, SLASH))
}

func stop(context *gin.Context, r TaskRequest, path string) {
	m := ProjectTasks(path, r.Project)
	ch, ok := m[r.Task]
	if ok {
		ch <- 1
		RemoveProjectTask(path, r.Project, r.Task)
	} else if r.Task == EMPTY {
		for _, ch := range m {
			ch <- 1
		}
		RemoveProject(path, r.Project)
	}
	context.String(http.StatusOK, EMPTY)
}

func RunTestReport(context *gin.Context) {
	task := context.Param(ParamTask)
	s := []string{context.Request.Host, WS, task}
	h := strings.Join(s, SLASH)
	data := `
<script>
var ws = new WebSocket("ws://` + h + `");   
ws.onopen = function(evt) {  
    console.log("Connection open ...");  
    ws.send("Hello WebSockets!");  
   document.open();
}; 
ws.onmessage = function(evt) {  
   document.write(evt.data);
};
//连接关闭时触发  
ws.onclose = function(evt) {  
    console.log("Connection closed.");  
   document.close();
}; 
</script>`
	context.Data(http.StatusOK, CTTEXT, []byte(data))
}

func CallPlugin(context *gin.Context) {
	plugin := context.Param(ParamPlugin)
	function := context.Param(ParamFunction)
	key := Join(COLON, plugin, function)
	bch := make(chan []byte)
	AnyChannel(ModuleDB) <- StoreType{Service: PROJECT, Op: Get, Key: key, BCh: bch}
	v := <-bch
	if len(v) == 0 {
		context.String(http.StatusOK, "No such plugin %s!", plugin)
		return
	}
	var q ControlInfo
	err := GobDecode(v, &q)
	if err != nil {
		context.String(http.StatusOK, "No such plugin %s!", plugin)
		return
	}
	context.JSON(http.StatusOK, q)
}

func setNoCache(context *gin.Context) {
	context.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Writer.Header().Set("Pragma", "no-cache")
	context.Writer.Header().Set("Expires", "0")
}
func triggerTask(ip, task string) {
	if !GetInit(ip) {
		SetInit(ip)
	}
	SetTask(task)
}

func handleWSTask(context *gin.Context, p []byte, task string) {
	ip := context.Request.Host
	switch {
	case MatchStart(p, WsInitTask):
		triggerTask(ip, task)
	case MatchStart(p, WsWaitTask):
		k := GetStatus(ip)
		for k != RunTask {
			time.Sleep(time.Duration(GetInterval()) * time.Millisecond)
			k = GetStatus(ip)
		}
		triggerTask(ip, task)
	case MatchStart(p, WsBeginTask):
		t := string(p)
		k := strings.Split(t, SPACE)
		Debug(ModuleControl, "begin TASK", k)
		triggerTask(ip, task)
		SetStatus(ip, RunTask)
	}
}

func RunWSTask(context *gin.Context) {
	task := context.Param(ParamTask)
	rch := BytesChannel(task) // 应答通道
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		Debug(ModuleControl, err)
		return
	}
	defer ws.Close()
	t, p, err := ws.ReadMessage()
	if err != nil {
		Debug(ModuleControl, err)
	}

	Debug(ModuleControl, t, string(p), "task", task)
	handleWSTask(context, p, task)
	Debug(ModuleControl, "sending task", task)
	for {
		if !GetTask(task) {
			break
		}
		//Debug(ModuleControl, "receiving task", task)
		b := <-rch
		//Debug(ModuleControl, "task", task, "sending", string(b))
		if len(b) == 0 {
			continue
		} else if bytes.Equal(b, []byte(EOF)) {
			break
		}
		err = ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			Debug(ModuleControl, err)
			break
		}
	}
	//RemoveTask(task)
}

func StopTask(context *gin.Context) {
	SetStatus(context.Request.Host, WaitTask)
	context.String(http.StatusOK, OK)
}

func Config(context *gin.Context) {
	key := context.Param(ParamKey)
	value := context.Param(ParamValue)
	k := 0
	var err error
	if key == TOTAL || key == INTERVAL {
		k, err = strconv.Atoi(value)
		if err != nil {
			context.String(http.StatusBadRequest, "%s的值只能是数字", key)
			return
		}
	}
	switch key {
	case LOG:
		SetLogLevel(value)
	case INTERVAL:
		SetInterval(k)
	default:
		context.String(http.StatusBadRequest, "请求类型不支持")
	}
	context.String(http.StatusOK, OK)
}
