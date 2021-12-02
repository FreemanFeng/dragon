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
