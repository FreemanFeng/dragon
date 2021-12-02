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

package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var mL *LogConfig
var oL sync.Once

func GetLogConfig() *LogConfig {
	oL.Do(func() {
		mL = &LogConfig{}
	})
	return mL
}

func InitLogConfig() {
	k := GetLogConfig()
	//k.LogLevel = DEBUG
	k.LogLevel = INFO
	k.LogALL = 1
}

func GetLogLevel() string {
	k := GetLogConfig()
	return k.LogLevel
}

func SetLogLevel(level string) {
	if level != INFO && level != DEBUG && level != DETAIL {
		return
	}
	k := GetLogConfig()
	k.LogLevel = level
}

func GetLogModule(module string) bool {
	k := GetLogConfig()
	if k.LogALL > 0 {
		return true
	}
	_, ok := k.LogModules[module]
	return ok
}

func SetLogModules(param string) {
	k := GetLogConfig()
	if len(param) == 0 {
		k.LogALL = 1
		return
	}
	k.LogModules = map[string]int{}
	modules := strings.Split(param, ",")
	for _, v := range modules {
		k.LogModules[v] = 1
	}
	if len(k.LogModules) > 0 {
		k.LogALL = 0
	}
}

func Now() string {
	return time.Now().Format(time.RFC850)
}

func Log2(params ...interface{}) {
	if len(params) == 0 {
		return
	}
	fmt.Print(Now(), " ")
	fmt.Println(params...)
}

func Log(params ...interface{}) {
	if len(params) == 0 {
		return
	}
	fmt.Print(Now(), " ", FuncName(), " ")
	fmt.Println(params...)
}

func Detail(params ...interface{}) {
	switch GetLogLevel() {
	case DETAIL:
		if !GetLogModule(params[0].(string)) {
			return
		}
		Log(params...)
	}
}

func Debug(params ...interface{}) {
	switch GetLogLevel() {
	case DETAIL:
		if !GetLogModule(params[0].(string)) {
			return
		}
		Log(params...)
	case DEBUG:
		if !GetLogModule(params[0].(string)) {
			return
		}
		Log(params...)
	}
}

func Info(params ...interface{}) {
	switch GetLogLevel() {
	case DETAIL:
		Log(params...)
	case DEBUG:
		Log(params...)
	case INFO:
		Log(params...)
	}
}

func Warning(params ...interface{}) {
	switch GetLogLevel() {
	case DETAIL:
		Log(params...)
	case DEBUG:
		Log(params...)
	case INFO:
		Log(params...)
	case WARNING:
		Log(params...)
	}

}

func Err(params ...interface{}) {
	switch GetLogLevel() {
	case DETAIL:
		Log(params...)
	case DEBUG:
		Log(params...)
	case INFO:
		Log(params...)
	case WARNING:
		Log(params...)
	case ERROR:
		Log(params...)
	}

}

func AssertOK(err error, params ...interface{}) {
	if err != nil {
		Log(params...)
		os.Exit(1)
	}
}

func FatalError(params ...interface{}) {
	name, file, line := Func(2)
	Log("Fatal Error in File", file, "Line", line, "Function", name)
	Log(params...)
}

func Func(skip int) (string, string, int) {
	fun, file, line, ok := runtime.Caller(skip)
	if !ok {
		Log("Could not Get runtime Caller info")
		os.Exit(1)
	}
	name := runtime.FuncForPC(fun).Name()
	return name, file, line
}

func FuncName() string {
	name, _, _ := Func(4)
	h := strings.Split(name, SLASH)
	if len(h) < 2 {
		return name
	}
	n := len(h)
	return h[n-1]
}
