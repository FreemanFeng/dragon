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
	"strings"
	"sync"
)

var mG *GlobalConfig
var oG sync.Once

func GetGlobalConfig() *GlobalConfig {
	oG.Do(func() {
		mG = &GlobalConfig{}
	})
	return mG
}

func InitGlobalConfig() {
	k := GetGlobalConfig()
	k.Paths = sync.Map{}
	k.Tasks = sync.Map{}
	k.Projects = sync.Map{}
	k.IPs = sync.Map{}
	k.Ports = sync.Map{}
	k.Calls = sync.Map{}
	k.Any = sync.Map{}
	k.Plugins = sync.Map{}
	SetIP(PubDNS, DefaultDNS)
}

func GetAny(key string, defaultValue interface{}) interface{} {
	k := GetGlobalConfig()
	t, ok := k.Any.Load(key)
	if !ok {
		k.Any.Store(key, defaultValue)
		return defaultValue
	}
	return t
}

func SetAny(key string, value interface{}) {
	k := GetGlobalConfig()
	k.Any.Store(key, value)
}

func GetCall(call string) interface{} {
	k := GetGlobalConfig()
	t, ok := k.Calls.Load(call)
	if !ok {
		Err(ModuleCommon, "无法找到内置函数", call)
		return nil
	}
	return t
}

func SetCall(call string, value interface{}) {
	k := GetGlobalConfig()
	k.Calls.Store(call, value)
}

func GetPath(key, defaultValue string) string {
	k := GetGlobalConfig()
	t, _ := k.Paths.LoadOrStore(key, defaultValue)
	return t.(string)
}

func SetPath(key, path string) {
	k := GetGlobalConfig()
	k.Paths.Store(key, path)
}

func SetIP(key, ip string) {
	k := GetGlobalConfig()
	k.IPs.Store(key, ip)
}

func GetIP(key string) string {
	k := GetGlobalConfig()
	t, ok := k.IPs.Load(key)
	if !ok {
		return EMPTY
	}
	return t.(string)
}

func SetPort(key string, port int) {
	k := GetGlobalConfig()
	k.Ports.Store(key, port)
}

func GetPort(key string) int {
	k := GetGlobalConfig()
	t, ok := k.Ports.Load(key)
	if !ok {
		return UNDEFINED
	}
	return t.(int)
}

func SetPlugin(service, plugin string) {
	k := GetGlobalConfig()
	k.Plugins.Store(service, plugin)
}

func GetPlugin(service, defaultValue string) string {
	k := GetGlobalConfig()
	t, _ := k.Plugins.LoadOrStore(service, defaultValue)
	return t.(string)
}

func SetStatus(ip string, status int) {
	k := GetGlobalConfig()
	k.Status.Store(ip, status)
}

func GetStatus(ip string) int {
	k := GetGlobalConfig()
	v, ok := k.Status.Load(ip)
	if !ok {
		return WaitTask
	}
	return v.(int)
}

func SetInit(ip string) {
	k := GetGlobalConfig()
	k.Init.Store(ip, true)
}

func GetInit(ip string) bool {
	k := GetGlobalConfig()
	_, ok := k.Init.Load(ip)
	return ok
}

func SetInterval(interval int) {
	k := GetGlobalConfig()
	k.Interval = interval
}

func GetInterval() int {
	k := GetGlobalConfig()
	return k.Interval
}

func SetTask(task string) {
	k := GetGlobalConfig()
	k.Tasks.Store(task, DONE)
}

func GetTask(task string) bool {
	k := GetGlobalConfig()
	_, ok := k.Tasks.Load(task)
	return ok
}
func RemoveTask(task string) {
	k := GetGlobalConfig()
	k.Tasks.Delete(task)
}

func ProjectTasks(path, key string) map[string]chan int {
	k := GetGlobalConfig()
	s := []string{path, key}
	h := strings.Join(s, COLON)
	v, _ := k.Projects.LoadOrStore(h, map[string]chan int{})
	return v.(map[string]chan int)
}
func SetProjectTask(path, key, task string, ch chan int) {
	k := GetGlobalConfig()
	s := []string{path, key}
	h := strings.Join(s, COLON)
	v := ProjectTasks(path, key)
	v[task] = ch
	k.Projects.Store(h, v)
}

func RemoveProject(path, key string) {
	k := GetGlobalConfig()
	s := []string{path, key}
	h := strings.Join(s, COLON)
	k.Projects.Delete(h)
}

func RemoveProjectTask(path, key, task string) {
	k := GetGlobalConfig()
	s := []string{path, key}
	h := strings.Join(s, COLON)
	v := ProjectTasks(path, key)
	delete(v, task)
	k.Projects.Store(h, v)
}
