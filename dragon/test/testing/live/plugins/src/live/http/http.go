package http

import (
	"fmt"
	"goplugin"
)

func OnReadySending(m, c, r map[string][]byte) {
	fmt.Println("OnReadySending message templates", m, "config", c, "request attributes", r)
	var b []byte
	ab := m["a"]
	fmt.Println("before update a", string(m["a"]))
	ma := goplugin.DecodeMap(ab)
	if ma == nil {
		return
	}
	ma["hello2"] = "world"
	b = goplugin.Encode(ma)
	if b == nil {
		return
	}
	m["a"] = goplugin.Encode(ma)
	fmt.Println("after update a", string(m["a"]))
}

func OnReceived(m, c, r map[string][]byte) {
	fmt.Println("OnReceived message templates", m, "config", c, "response", r)
	mb := r["resp"]
	fmt.Println("before update resp", string(r["resp"]))
	mm := goplugin.DecodeMap(mb)
	if mm == nil {
		return
	}
	mm["hello2"] = "world"
	b := goplugin.Encode(mm)
	if b == nil {
		return
	}
	r["resp"] = b
	fmt.Println("after update resp", string(r["resp"]))
}

func OnError(m, c, r map[string][]byte) {
	fmt.Println("OnError message templates", m, "config", c, "response", r)
	k, ok := r["code"]
	if ok {
		code := string(k)
		if code == "400" {
			r["code"] = []byte("500")
		}
	}
}
