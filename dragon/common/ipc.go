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
	"sync"
)

var mC *IPCType
var oC sync.Once

func GetIPCInstance() *IPCType {
	oC.Do(func() {
		mC = &IPCType{}
	})
	return mC
}

type IPCType struct {
	// 协程间通信管道
	Any       sync.Map
	Int       sync.Map
	Byte      sync.Map
	String    sync.Map
	Bytes     sync.Map
	AnyMap    sync.Map
	IntMap    sync.Map
	StringMap sync.Map
	BytesMap  sync.Map
}

func InitChannels() {
	k := GetIPCInstance()
	// 协程间通信管道
	k.Any = sync.Map{}
	k.Int = sync.Map{}
	k.Byte = sync.Map{}
	k.String = sync.Map{}
	k.Bytes = sync.Map{}
	k.AnyMap = sync.Map{}
	k.IntMap = sync.Map{}
	k.StringMap = sync.Map{}
	k.BytesMap = sync.Map{}
}

func AnyChannel(key interface{}) chan interface{} {
	k := GetIPCInstance()
	v, _ := k.Any.LoadOrStore(key, make(chan interface{}))
	return v.(chan interface{})
}

func IntChannel(key interface{}) chan int {
	k := GetIPCInstance()
	v, _ := k.Int.LoadOrStore(key, make(chan int))
	return v.(chan int)
}

func StrChannel(key interface{}) chan string {
	k := GetIPCInstance()
	v, _ := k.String.LoadOrStore(key, make(chan string))
	return v.(chan string)
}

func ByteChannel(key interface{}) chan byte {
	k := GetIPCInstance()
	v, _ := k.Byte.LoadOrStore(key, make(chan byte))
	return v.(chan byte)
}

func BytesChannel(key interface{}) chan []byte {
	k := GetIPCInstance()
	v, _ := k.Bytes.LoadOrStore(key, make(chan []byte))
	return v.(chan []byte)
}

func AnyMapChannel(key interface{}) chan map[interface{}]interface{} {
	k := GetIPCInstance()
	v, _ := k.AnyMap.LoadOrStore(key, make(chan map[interface{}]interface{}))
	return v.(chan map[interface{}]interface{})
}

func IntMapChannel(key interface{}) chan map[string]int {
	k := GetIPCInstance()
	v, _ := k.IntMap.LoadOrStore(key, make(chan map[string]int))
	return v.(chan map[string]int)
}

func StrMapChannel(key interface{}) chan map[string]string {
	k := GetIPCInstance()
	v, _ := k.StringMap.LoadOrStore(key, make(chan map[string]string))
	return v.(chan map[string]string)
}

func BytesMapChannel(key interface{}) chan map[string][]byte {
	k := GetIPCInstance()
	v, _ := k.BytesMap.LoadOrStore(key, make(chan map[string][]byte))
	return v.(chan map[string][]byte)
}
