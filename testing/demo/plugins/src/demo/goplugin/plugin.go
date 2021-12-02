package goplugin

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var mG *Config
var oG sync.Once

func GetConfig() *Config {
	oG.Do(func() {
		mG = &Config{Calls: sync.Map{}}
	})
	return mG
}

func Func(skip int) (string, string, int) {
	fun, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("Could not Get runtime Caller info")
		os.Exit(1)
	}
	name := runtime.FuncForPC(fun).Name()
	return name, file, line
}

func FuncName(skip int) string {
	name, _, _ := Func(skip)
	h := strings.Split(name, SLASH)
	if len(h) < 2 {
		return name
	}
	n := len(h)
	return h[n-1]
}

func FolderName(skip int) string {
	_, file, _ := Func(skip)
	h := strings.Split(file, SLASH)
	if len(h) < 3 {
		return file
	}
	n := len(h)
	return h[n-2]
}

func GetCall(call string) reflect.Value {
	k := GetConfig()
	s := []string{FolderName(5), call}
	key := strings.Join(s, DOT)
	f, ok := k.Calls.Load(key)
	fmt.Println("Get Call", key, "to", f)
	if !ok {
		fmt.Println("could not get func", call)
		return reflect.ValueOf(nil)
	}
	return reflect.ValueOf(f)
}

func GetXCall(call string) reflect.Value {
	k := GetConfig()
	s := []string{k.Project, call}
	key := strings.Join(s, DOT)
	f, ok := k.Calls.Load(key)
	fmt.Println("Get Call", key, "to", f)
	if !ok {
		fmt.Println("could not get func", call)
		return reflect.ValueOf(nil)
	}
	return reflect.ValueOf(f)
}

func SetCall(call string, value interface{}) {
	k := GetConfig()
	s := []string{FolderName(4), call}
	key := strings.Join(s, DOT)
	fmt.Println("Set Call", key, "to", value)
	k.Calls.Store(key, value)
	k.Project = FolderName(4)
}

func Init(m map[string]interface{}) {
	var s []string
	for k, v := range m {
		SetCall(k, v)
		s = append(s, k)
	}
}

func Run(name string, b []byte) []byte {
	var bs [][]byte
	e := json.Unmarshal(b, &bs)
	if e != nil {
		fmt.Println("plugin func", name, "deserialized failed", e)
		return nil
	}
	return Process(name, bs)
}

func Process(name string, bs [][]byte) []byte {
	var e error
	fmt.Println("run plugin func", name)
	n := len(bs) / 2
	ms := make([]map[string][]byte, n)
	rs := make([][]byte, n)
	in := make([]reflect.Value, n)
	f := GetCall(name)
	if f == reflect.ValueOf(nil) {
		fmt.Println(">>>> Failed in calling function", name)
		return []byte(EMPTY)
	}
	for i := 0; i < n; i++ {
		k := bs[i*2]
		v := bs[i*2+1]
		switch string(k) {
		case "MAP":
			e = json.Unmarshal(v, &ms[i])
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument deserialized failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(ms[i])
		}
	}
	if strings.Contains(name, ".") {
		f.Call(in)
		for i := range in {
			rs[i], e = json.Marshal(in[i].Interface())
			if e != nil {
				fmt.Println("plugin func", name, "result serialized failed", e)
				return nil
			}
		}
		b, e := json.Marshal(rs)
		if e != nil {
			fmt.Println("plugin func", name, "result serialized failed", e)
			return nil
		}
		return b
	}

	for i := 0; i < n; i++ {
		k := bs[i*2]
		v := bs[i*2+1]
		switch string(k) {
		case "MAP":
			var m map[string]interface{}
			e = json.Unmarshal(v, &m)
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument deserialized failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(m)
		case "LIST":
			var m []interface{}
			e = json.Unmarshal(v, &m)
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument deserialized failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(m)
		case "STR":
			in[i] = reflect.ValueOf(string(v))
		case "BYTES":
			in[i] = reflect.ValueOf(v)
		case "INT":
			x, e := strconv.ParseInt(string(v), 10, 64)
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument convert to INT type failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(x)
		case "UINT":
			x, e := strconv.ParseUint(string(v), 10, 64)
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument convert to UINT type failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(x)
		case "FLOAT":
			x, e := strconv.ParseFloat(string(v), 10)
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument convert to FLOAT type failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(x)
		case "BOOL":
			x, e := strconv.ParseBool(string(v))
			if e != nil {
				fmt.Println("plugin func", name, "No.", i, "argument convert to BOOL type failed", e)
				return nil
			}
			in[i] = reflect.ValueOf(x)
		default:
			in[i] = reflect.ValueOf(nil)
		}
	}
	r := f.Call(in)
	if len(r) == 0 {
		return []byte(EMPTY)
	}
	k := r[0].Interface()
	return k.([]byte)
}

func DecodeMap(b []byte) map[string]interface{} {
	var v map[string]interface{}
	e := json.Unmarshal(b, &v)
	if e != nil {
		fmt.Println("JSON deserialized failed", e)
		return nil
	}
	return v
}

func DecodeList(b []byte) []interface{} {
	var v []interface{}
	e := json.Unmarshal(b, &v)
	if e != nil {
		fmt.Println("JSON deserialized failed", e)
		return nil
	}
	return v
}

func Encode(v interface{}) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		fmt.Println("JSON serialized failed", e)
		return nil
	}
	return b
}
