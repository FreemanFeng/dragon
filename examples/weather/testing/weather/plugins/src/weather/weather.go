package main

import (
	"weather/common"
	"weather/http"
)

// 插件初始化，返回回调函数列表
func Init() map[string]func(params ...interface{}) interface{} {
	m := map[string]func(params ...interface{}) interface{}{}
	m["http.OnStart"] = http.OnStart
	m["http.ToSend"] = http.ToSend
	m["http.ToAuth"] = http.ToAuth
	m["http.OnConnected"] = http.OnConnected
	m["http.OnSending"] = http.OnSending
	m["http.OnAuthorizing"] = http.OnAuthorizing
	m["http.OnSent"] = http.OnSent
	m["http.OnAuthorized"] = http.OnAuthorized
	m["http.ToReceive"] = http.ToReceive
	m["http.OnReceiving"] = http.OnReceiving
	m["http.OnReceived"] = http.OnReceived
	m["http.OnError"] = http.OnError
	m["http.OnDisconnected"] = http.OnDisconnected
	m["http.OnTimeout"] = http.OnTimeout
	m["http.OnEnd"] = http.OnEnd
	m["Sign"] = common.Sign
	return m
}

func Test() {
	Init()
}
