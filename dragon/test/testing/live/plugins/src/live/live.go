package main

import "C"

import (
	"demo/common"
	"demo/http"
	"demo/pb"
	"flag"
	"fmt"
	"os"

	"goplugin"
)

//export Init
func Init() map[string]interface{} {
	s := map[string]interface{}{
		"http.OnReadySending": http.OnReadySending,
		"http.OnReceived":     http.OnReceived,
		"http.OnError":        http.OnError,
		"pb.OnReadySending":   pb.OnReadySending,
		"pb.OnReceived":       pb.OnReceived,
		"pb.OnError":          pb.OnError,
		"Sign":                common.Sign}
	goplugin.Init(s)
	return s
}

//export Run
func Run(name string, b []byte) []byte {
	return goplugin.Run(name, b)
}

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], "[-p port][-h][-d]")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var help bool
	var port int
	flag.BoolVar(&help, "h", false, "Show Usage")
	flag.IntVar(&port, "p", 8088, "Control Port")
	flag.Parse()
	if help {
		Usage()
		return
	}
	goplugin.Serve(port, Init, Run)
}
