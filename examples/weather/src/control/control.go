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
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/FreemanFeng/dragon/examples/weather/src/common"

	"github.com/FreemanFeng/dragon/examples/weather/src/services"

	"github.com/gin-gonic/gin"
)

func Run(port int) {
	rand.Seed(time.Now().Unix())
	// 初始化引擎
	engine := gin.Default()
	// 注册路由和处理函数
	engine.Any(PathRoot, Home)
	engine.POST(PathGetWeather, GetWeather)

	// 绑定端口，然后启动应用
	s := []string{COLON, strconv.Itoa(port)}
	engine.Run(strings.Join(s, EMPTY))
}

func GetWeather(context *gin.Context) {
	s := context.Query(ParamTS)
	ts, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		fmt.Println("ts解析出错", e)
		context.String(http.StatusBadRequest, EMPTY)
		return
	}
	s = context.Query(ParamNonce)
	nonce, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		fmt.Println("nonce解析出错", e)
		context.String(http.StatusBadRequest, EMPTY)
		return
	}
	sign := context.Query(ParamSign)
	msg := &WeatherRequest{}
	e = context.BindJSON(msg)
	if e != nil {
		fmt.Println("msg解析出错", e)
		context.String(http.StatusBadRequest, EMPTY)
		return
	}
	if msg.Interval != "小时" && msg.Interval != "天" {
		context.String(http.StatusBadRequest, "请求字段interval只能是小时或天")
		return
	}
	if msg.City == "" {
		context.String(http.StatusBadRequest, "请求字段city为空")
		return
	}
	if services.ExistsNonce(nonce) {
		fmt.Println("Nonce", nonce, "重复")
		//context.String(http.StatusBadRequest, "Nonce", nonce, "重复")
		//return
	}
	fmt.Println("收到请求:", msg)
	data := services.GetWeather(ts, nonce, sign, msg)
	context.JSON(http.StatusOK, data)
}

func Home(context *gin.Context) {
	context.String(http.StatusOK, OK)
}
