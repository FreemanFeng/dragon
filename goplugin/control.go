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

package goplugin

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(port int, initFunc, runFunc interface{}) {
	rand.Seed(time.Now().Unix())
	engine := gin.Default()
	SetCall(INIT, initFunc)
	SetCall(RUN, runFunc)
	engine.Any(PathRoot, Home)
	engine.GET(PathInit, InitService)
	engine.POST(PathRunService, RunService)

	s := []string{COLON, strconv.Itoa(port)}
	engine.Run(strings.Join(s, EMPTY))
}

func InitService(context *gin.Context) {
	var s []string
	x := GetXCall(INIT)
	f := x.Interface().(func() map[string]interface{})
	m := f()
	for k := range m {
		s = append(s, k)
	}
	fmt.Println(">>>> 已初始化RPC服务列表:", s)
	context.JSON(http.StatusOK, s)
}

func RunService(context *gin.Context) {
	s := context.Param(ParamService)
	r, e := context.GetRawData()
	if e != nil {
		fmt.Println("RunService failed", e)
		context.String(http.StatusBadRequest, e.Error())
		return
	}
	x := GetXCall(RUN)
	f := x.Interface().(func(name string, b []byte) []byte)
	b := f(s, r)
	context.Data(http.StatusOK, CTJSON, b)
}

func Home(context *gin.Context) {
	context.String(http.StatusOK, OK)
}
