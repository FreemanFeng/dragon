package util

import (
	"reflect"
)

func Call(funcName interface{}, params []interface{}) []reflect.Value {
	f := reflect.ValueOf(funcName)
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in)
}

func CallFunc(funcName interface{}, params ...interface{}) []reflect.Value {
	return Call(funcName, params)
}

func CallRefs(funcName interface{}, params []interface{}) []reflect.Value {
	refs := make([]interface{}, len(params))
	for k := range params {
		refs[k] = &params[k]
	}
	f := reflect.ValueOf(funcName)
	in := make([]reflect.Value, len(refs))
	for k, param := range refs {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in)
}
