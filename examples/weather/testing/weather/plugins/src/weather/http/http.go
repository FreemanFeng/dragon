package http

import "fmt"

func OnStart(params ...interface{}) interface{} {
	fmt.Print("On Start")
	return ""
}

func ToSend(params ...interface{}) interface{} {
	fmt.Print("To Send")
	return ""
}

func ToAuth(params ...interface{}) interface{} {
	fmt.Print("To Auth")
	return ""
}

func OnConnected(params ...interface{}) interface{} {
	fmt.Print("On Connected")
	return ""
}

func OnSending(params ...interface{}) interface{} {
	fmt.Print("On Sending")
	return ""
}

func OnAuthorizing(params ...interface{}) interface{} {
	fmt.Print("On Authorizing")
	return ""
}

func OnSent(params ...interface{}) interface{} {
	fmt.Print("On Sent")
	return ""
}

func OnAuthorized(params ...interface{}) interface{} {
	fmt.Print("On Authorized")
	return ""
}

func ToReceive(params ...interface{}) interface{} {
	fmt.Print("To Receive")
	return ""
}

func OnReceiving(params ...interface{}) interface{} {
	fmt.Print("On Receiving")
	return ""
}

func OnReceived(params ...interface{}) interface{} {
	fmt.Print("On Received")
	return ""
}

func OnError(params ...interface{}) interface{} {
	fmt.Print("On Error")
	return ""
}

func OnDisconnected(params ...interface{}) interface{} {
	fmt.Print("On Disconnected")
	return ""
}

func OnTimeout(params ...interface{}) interface{} {
	fmt.Print("On Timeout")
	return ""
}

func OnEnd(params ...interface{}) interface{} {
	fmt.Print("On End")
	return ""
}
