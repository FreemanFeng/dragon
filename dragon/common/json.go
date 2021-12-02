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
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func WalkJSON(buf []byte, key string) (interface{}, map[string][]FormatType) {
	var msg interface{}
	fields := map[string][]FormatType{}
	markBrackets := []byte(LeftBracket)
	markBrace := []byte(LeftBrace)
	b := bytes.TrimSpace(buf)
	if bytes.Index(b, markBrackets) == 0 {
		msg = ParseJsonList(b, key, fields, []FormatType{})
	} else if bytes.Index(b, markBrace) == 0 {
		msg = ParseJsonMap(b, key, fields, []FormatType{})
	}
	return msg, fields
}

func ToJSON(buf []byte) interface{} {
	var msg interface{}
	if buf == nil {
		return nil
	}
	markBrackets := []byte(LeftBracket)
	markBrace := []byte(LeftBrace)
	b := bytes.TrimSpace(buf)
	if bytes.Index(b, markBrackets) == 0 {
		msg = ToJsonList(b)
	} else if bytes.Index(b, markBrace) == 0 {
		msg = ToJsonMap(b)
	} else {
		return nil
	}
	return msg
}

func IsJsonList(buf []byte) bool {
	var msg []interface{}
	b := bytes.TrimSpace(buf)
	err := json.Unmarshal(b, &msg)
	if err != nil {
		return false
	}
	return true
}

func IsJsonMap(buf []byte) bool {
	var msg map[string]interface{}
	b := bytes.TrimSpace(buf)
	err := json.Unmarshal(b, &msg)
	if err != nil {
		return false
	}
	return true
}

func ParseJsonList(buf []byte, key string, fields map[string][]FormatType, keys []FormatType) []interface{} {
	var msg []interface{}
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		Err(ModuleCommon, "Failed in Parsing Json Array Data", err)
	}
	WalkListFields(msg, key, fields, keys)
	return msg
}

func ToJsonList(buf []byte) []interface{} {
	var msg []interface{}
	//err := json.Unmarshal(buf, &msg)
	err := Unmarshal(buf, &msg)
	if err != nil {
		Err(ModuleCommon, "!!!!!!Failed in Converting Json Array Data", err)
		return nil
	}
	return msg
}

func WalkListFields(data []interface{}, key string, fields map[string][]FormatType, keys []FormatType) {
	for i, v := range data {
		t := Join(EMPTY, key, LeftBracket, strconv.Itoa(i), RightBracket)
		k := FormatType{Content: strconv.Itoa(i), IsInt: true}
		WalkFields(v, t, fields, append(keys, k))
	}
}

func ParseJsonMap(buf []byte, key string, fields map[string][]FormatType, keys []FormatType) map[string]interface{} {
	var msg map[string]interface{}
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		Err(ModuleCommon, "Failed in Parsing Json Map Data", err)
	}
	WalkMapFields(msg, key, fields, keys)
	return msg
}

func ToJsonMap(buf []byte) map[string]interface{} {
	var msg map[string]interface{}
	//err := json.Unmarshal(buf, &msg)
	err := Unmarshal(buf, &msg)
	if err != nil {
		Err(ModuleCommon, "!!!!!!Failed in Converting Json Map Data", err)
		return nil
	}
	return msg
}

func WalkMapFields(data map[string]interface{}, key string, fields map[string][]FormatType, keys []FormatType) {
	for k, v := range data {
		t := Join(DOT, key, k)
		s := FormatType{Content: k}
		WalkFields(v, t, fields, append(keys, s))
	}
}

func WalkFields(data interface{}, key string, fields map[string][]FormatType, keys []FormatType) {
	if IsList(data) {
		WalkListFields(data.([]interface{}), key, fields, keys)
	} else if IsMap(data) {
		WalkMapFields(data.(map[string]interface{}), key, fields, keys)
	} else {
		fields[key] = keys
	}
}

func UpdateJsonOnKeys(msg interface{}, value []interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return msg
	}
	top := UNDEFINED
	if keys[0].IsInt {
		if m, ok := msg.([]interface{}); ok {
			for top < len(value) {
				top, msg = UpdateJsonListOnKeys(m, value, keys, top+1)
			}
			return msg
		}
	}
	for top < len(value) {
		top, msg = UpdateJsonMapOnKeys(msg.(map[string]interface{}), value, keys, top+1)
	}
	return msg
}

func UpdateJsonListOnKeys(data []interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	// 没有字段可搜索，直接返回
	if len(keys) == 0 || top >= len(value) {
		return top, data
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return top, data
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for i := range data {
				if keys[1].IsInt {
					k, _ := strconv.Atoi(keys[1].Content)
					if k == i {
						return UpdateJsonListOnKeys(data, value, keys[1:], top)
					}
				}
				top, data[i] = UpdateFieldsOnKeys(data[i], value, keys, top)
				if top == len(value) {
					return top, data
				}
			}
			return top, data
		}
		return UpdateJsonListOnKeys(data, value, keys[1:], top)
	}
	i, _ := strconv.Atoi(keys[0].Content)
	n := len(data)
	if i >= n || top >= len(value) {
		return top, data
	}
	if len(keys) == 1 {
		data[i] = AlignType(data[i], value[top])
		return top + 1, data
	}
	top, data[i] = UpdateFieldsOnKeys(data[i], value, keys[1:], top)
	return top, data
}

func UpdateJsonMapOnKeys(data map[string]interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	// 没有字段可搜索，直接返回
	if len(keys) == 0 || top >= len(value) {
		return top, data
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return top, data
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for k := range data {
				if k == keys[1].Content {
					return UpdateJsonMapOnKeys(data, value, keys[1:], top)
				}
				top, data[k] = UpdateFieldsOnKeys(data[k], value, keys, top)
				if top == len(value) {
					return top, data
				}
			}
			return top, data
		}
		return UpdateJsonMapOnKeys(data, value, keys[1:], top)
	}
	k := keys[0].Content
	if _, ok := data[k]; !ok {
		return top, data
	}
	if len(keys) == 1 {
		data[k] = AlignType(data[k], value[top])
		return top + 1, data
	}
	top, data[k] = UpdateFieldsOnKeys(data[k], value, keys[1:], top)
	return top, data
}

func UpdateFieldsOnKeys(data interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	if IsList(data) {
		top, data = UpdateJsonListOnKeys(data.([]interface{}), value, keys, top)
	} else if IsMap(data) {
		top, data = UpdateJsonMapOnKeys(data.(map[string]interface{}), value, keys, top)
	}
	return top, data
}

func UpdateJSON(msg interface{}, value interface{}) interface{} {
	if msg == nil {
		return value
	}
	if IsList(msg) && IsList(value) {
		msg = UpdateJsonList(msg.([]interface{}), value.([]interface{}))
	}
	if IsMap(msg) && IsMap(value) {
		msg = UpdateJsonMap(msg.(map[string]interface{}), value.(map[string]interface{}))
	}
	return msg
}

func UpdateJsonList(data []interface{}, value []interface{}) interface{} {
	n := len(value)
	if n == 0 {
		return data
	}
	for i := 0; i < n; i++ {
		if i >= len(data) {
			data = append(data, value[i])
			continue
		}
		if !IsList(data[i]) && !IsList(value[i]) && !IsMap(data[i]) && !IsMap(value[i]) {
			data[i] = value[i]
			continue
		}
		data[i] = UpdateFields(data[i], value[i])
	}
	return data
}

func UpdateJsonMap(data map[string]interface{}, value map[string]interface{}) interface{} {
	n := len(value)
	if n == 0 {
		return data
	}
	for k, v := range value {
		_, ok := data[k]
		if !ok {
			data[k] = v
			continue
		}
		if !IsList(data[k]) && !IsList(v) && !IsMap(data[k]) && !IsMap(v) {
			data[k] = v
		}
		data[k] = UpdateFields(data[k], v)
	}
	return data
}

func UpdateFields(data interface{}, value interface{}) interface{} {
	if IsList(data) && IsList(value) {
		data = UpdateJsonList(data.([]interface{}), value.([]interface{}))
	}
	if IsMap(data) && IsMap(value) {
		data = UpdateJsonMap(data.(map[string]interface{}), value.(map[string]interface{}))
	}
	return data
}

func AddJSON(msg interface{}, value []interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return msg
	}
	top := UNDEFINED
	if keys[0].IsInt {
		if m, ok := msg.([]interface{}); ok {
			for top < len(value) {
				top, msg = AddJsonList(m, value, keys, top+1)
			}
			return msg
		}
	}
	for top < len(value) {
		top, msg = AddJsonMap(msg.(map[string]interface{}), value, keys, top+1)
	}
	return msg
}

func AddJsonList(data []interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	// 没有字段可搜索，直接返回
	if len(keys) == 0 || top >= len(value) {
		return top, data
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return top, data
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for i := range data {
				if keys[1].IsInt {
					k, _ := strconv.Atoi(keys[1].Content)
					if k == i {
						return AddJsonList(data, value, keys[1:], top)
					}
				}
				top, data[i] = AddFields(data[i], value, keys, top)
				if top == len(value) {
					return top, data
				}
			}
			return top, data
		}
		return AddJsonList(data, value, keys[1:], top)
	}
	i, _ := strconv.Atoi(keys[0].Content)
	n := len(data)
	for n < i+1 {
		data = append(data, data[n-1])
		n++
	}
	if len(keys) == 1 {
		data[i] = value[top]
		return top + 1, data
	}
	top, data[i] = AddFields(data[i], value, keys[1:], top)
	return top, data
}

func AddJsonMap(data map[string]interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	// 没有字段可搜索，直接返回
	if len(keys) == 0 || top >= len(value) {
		return top, data
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return top, data
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for k := range data {
				if k == keys[1].Content {
					return AddJsonMap(data, value, keys[1:], top)
				}
				top, data[k] = AddFields(data[k], value, keys, top)
				if top == len(value) {
					return top, data
				}
			}
			return top, data
		}
		return AddJsonMap(data, value, keys[1:], top)
	}
	k := keys[0].Content
	if _, ok := data[k]; !ok {
		if len(keys) > 1 && keys[1].Content != EMPTY {
			if keys[1].IsInt {
				data[k] = []interface{}{}
			} else {
				data[k] = map[string]interface{}{}
			}
		}
	}
	if len(keys) == 1 {
		data[k] = value[top]
		return top + 1, data
	}
	top, data[k] = AddFields(data[k], value, keys[1:], top)
	return top, data
}

func AddFields(data interface{}, value []interface{}, keys []FormatType, top int) (int, interface{}) {
	if IsList(data) {
		top, data = AddJsonList(data.([]interface{}), value, keys, top)
	} else if IsMap(data) {
		top, data = AddJsonMap(data.(map[string]interface{}), value, keys, top)
	}
	return top, data
}

func DeleteJSON(msg interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return msg
	}
	Info(ModuleCommon, "消息", msg, "keys:", keys)
	if keys[0].IsInt {
		if m, ok := msg.([]interface{}); ok {
			return DeleteJsonList(m, keys)
		}
	}
	return DeleteJsonMap(msg.(map[string]interface{}), keys)
}

func DeleteJsonList(data []interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if keys[0].Content == EMPTY {
		//Log("---------------- 1")
		// 无效
		if len(keys) == 1 {
			//Log("---------------- 2")
			return nil
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			//Log("---------------- 3")

			for i := range data {
				//Log("---------------- 4")

				if keys[1].IsInt {
					//Log("---------------- 5")
					k, _ := strconv.Atoi(keys[1].Content)
					if k == i {
						//Log("---------------- 6")
						return DeleteJsonList(data, keys[1:])
					}
				}
				x := DeleteFields(data[i], keys)
				if x != nil {
					//Log("---------------- 7")
					return x
				}
			}
			return nil
		}
		//Log("---------------- 8")
		return DeleteJsonList(data, keys[1:])
	}
	//Log("---------------- 9")
	i, _ := strconv.Atoi(keys[0].Content)
	n := len(data)
	if i >= n {
		//Log("---------------- 10")
		return data
	}
	//Log("---------------- 11")
	if len(keys) == 1 {
		//Log("---------------- 12")
		if i == n-1 {
			//Log("---------------- 13")
			x := data[:i]
			//Info(ModuleCommon, "删除后消息内容为", x)
			return x
		}
		//Log("---------------- 14")
		x := append(data[:i], data[i+1:]...)
		//Info(ModuleCommon, "删除后消息内容为", x)
		return x
	}
	//Log("---------------- 15")
	data[i] = DeleteFields(data[i], keys[1:])
	return data
}

func DeleteJsonMap(data map[string]interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if keys[0].Content == EMPTY {
		//Log("++++++++++++ 1")
		// 无效
		if len(keys) == 1 {
			return nil
		}
		//Log("++++++++++++ 2")

		if len(keys) > 1 && keys[1].Content != EMPTY {
			//Log("++++++++++++ 3", keys)

			for k := range data {
				//Log("++++++++++++ 4", k, data)
				if k == keys[1].Content {
					//Log("++++++++++++ 5", data)
					return DeleteJsonMap(data, keys[1:])
				}
				//Log("++++++++++++ 6", k, data[k])
				x := DeleteFields(data[k], keys)
				if x != nil {
					//Log("++++++++++++ 7", x)
					return x
				}
			}
			return nil
		}
		//Log("++++++++++++ 8")
		return DeleteJsonMap(data, keys[1:])
	}
	//Log("++++++++++++ 9")

	k := keys[0].Content
	if len(keys) == 1 {
		//Log("++++++++++++ 10")
		//Info(ModuleCommon, "删除前字典消息内容1", data)
		delete(data, k)
		//Info(ModuleCommon, "删除后字典消息内容1", data)
		return data
	}
	//Log("++++++++++++ 11")
	//Info(ModuleCommon, "删除前字典消息内容2", data)
	x := DeleteFields(data[k], keys[1:])
	if x != nil {
		data = x.(map[string]interface{})
	}
	//Info(ModuleCommon, "删除后字典消息内容2", data)
	return data
}

func DeleteFields(data interface{}, keys []FormatType) interface{} {
	if IsList(data) {
		data = DeleteJsonList(data.([]interface{}), keys)
	} else if IsMap(data) {
		data = DeleteJsonMap(data.(map[string]interface{}), keys)
	}
	return nil
}

func IsField(msg interface{}, keys []FormatType) bool {
	if len(keys) == 0 {
		return false
	}
	switch keys[0].IsInt {
	case true:
		if m, ok := msg.([]interface{}); ok {
			return CheckJsonList(m, keys)
		}
	}
	return CheckJsonMap(msg.(map[string]interface{}), keys)
}

func CheckJsonList(data []interface{}, keys []FormatType) bool {
	i, _ := strconv.Atoi(keys[0].Content)
	n := len(data)
	if i >= n {
		return false
	}
	if len(keys) == 1 {
		return true
	}
	return CheckFields(data[i], keys[1:])
}

func CheckJsonMap(data map[string]interface{}, keys []FormatType) bool {
	k := keys[0].Content
	if _, ok := data[k]; !ok {
		return false
	}
	if len(keys) == 1 {
		return true
	}
	return CheckFields(data[k], keys[1:])
}

func CheckFields(data interface{}, keys []FormatType) bool {
	if IsList(data) {
		return CheckJsonList(data.([]interface{}), keys)
	} else if IsMap(data) {
		return CheckJsonMap(data.(map[string]interface{}), keys)
	}
	return false
}

func GetJSON(msg interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if keys[0].IsInt {
		if m, ok := msg.([]interface{}); ok {
			return GetFromJsonList(m, keys)
		}
	}
	return GetFromJsonMap(msg.(map[string]interface{}), keys)
}

func GetFromJsonList(data []interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return nil
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for i := range data {
				if keys[1].IsInt {
					k, _ := strconv.Atoi(keys[1].Content)
					if k == i {
						return GetFromJsonList(data, keys[1:])
					}
				}
				x := GetFromFields(data[i], keys)
				if x != nil {
					return x
				}
			}
			return nil
		}
		return GetFromJsonList(data, keys[1:])
	}

	n := len(data)
	i, e := strconv.Atoi(keys[0].Content)
	if e != nil {
		if Equal(keys[0].Content, RATotal, RATotalZ) {
			return n
		}
		return nil
	}
	if i >= n {
		return nil
	}
	if len(keys) == 1 {
		return data[i]
	}
	return GetFromFields(data[i], keys[1:])
}

func GetFromJsonMap(data map[string]interface{}, keys []FormatType) interface{} {
	if len(keys) == 0 {
		return nil
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return nil
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for k := range data {
				if k == keys[1].Content {
					return GetFromJsonMap(data, keys[1:])
				}
				x := GetFromFields(data[k], keys)
				if x != nil {
					return x
				}
			}
			return nil
		}
		return GetFromJsonMap(data, keys[1:])
	}
	k := keys[0].Content
	if _, ok := data[k]; !ok {
		if Equal(k, RATotal, RATotalZ) {
			return len(data)
		}
		return nil
	}
	if len(keys) == 1 {
		return data[k]
	}
	return GetFromFields(data[k], keys[1:])
}

func GetFromFields(data interface{}, keys []FormatType) interface{} {
	if IsList(data) {
		return GetFromJsonList(data.([]interface{}), keys)
	} else if IsMap(data) {
		return GetFromJsonMap(data.(map[string]interface{}), keys)
	}
	return nil
}

func GetAllJSON(msg interface{}, keys []FormatType, value map[string]interface{}) {
	if len(keys) == 0 {
		return
	}
	if keys[0].IsInt {
		if m, ok := msg.([]interface{}); ok {
			GetAllFromJsonList(m, keys, value, EMPTY)
			return
		}
	}
	GetAllFromJsonMap(msg.(map[string]interface{}), keys, value, EMPTY)
}

func GetAllFromJsonList(data []interface{}, keys []FormatType, value map[string]interface{}, path string) {
	if len(keys) == 0 {
		return
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for i := range data {
				s := Join(EMPTY, path, LeftBracket, strconv.Itoa(i), RightBracket)
				if keys[1].IsInt {
					k, _ := strconv.Atoi(keys[1].Content)
					if k == i {
						GetAllFromJsonList(data, keys[1:], value, s)
					}
				}
				GetAllFromFields(data[i], keys, value, s)
			}
			return
		}
		GetAllFromJsonList(data, keys[1:], value, path)
		return
	}

	n := len(data)
	i, e := strconv.Atoi(keys[0].Content)
	if e != nil {
		if Equal(keys[0].Content, RATotal, RATotalZ) {
			value[path] = n
		}
		return
	}
	if i >= n {
		return
	}
	if len(keys) == 1 {
		value[path] = data[i]
		return
	}
	GetAllFromFields(data[i], keys[1:], value, path)
}

func GetAllFromJsonMap(data map[string]interface{}, keys []FormatType, value map[string]interface{}, path string) {
	if len(keys) == 0 {
		return
	}
	if keys[0].Content == EMPTY {
		// 无效
		if len(keys) == 1 {
			return
		}
		if len(keys) > 1 && keys[1].Content != EMPTY {
			for k := range data {
				s := Join(DOT, path, k)
				if k == keys[1].Content {
					GetAllFromJsonMap(data, keys[1:], value, path)
				}
				GetAllFromFields(data[k], keys, value, s)
			}
			return
		}
		GetAllFromJsonMap(data, keys[1:], value, path)
	}
	k := keys[0].Content
	if _, ok := data[k]; !ok {
		if Equal(k, RATotal, RATotalZ) {
			value[path] = len(data)
		}
		return
	}
	if len(keys) == 1 {
		if strings.Index(path, DOT) == 0 {
			x := strings.Split(path, DOT)
			xn := len(x)
			path = strings.Join(x[1:xn], DOT)
		}
		value[path] = data[k]
		return
	}
	s := Join(DOT, path, k)
	if path == EMPTY {
		s = k
	}
	GetAllFromFields(data[k], keys[1:], value, s)
}

func GetAllFromFields(data interface{}, keys []FormatType, value map[string]interface{}, path string) {
	if IsList(data) {
		GetAllFromJsonList(data.([]interface{}), keys, value, path)
	} else if IsMap(data) {
		GetAllFromJsonMap(data.(map[string]interface{}), keys, value, path)
	}
}

func Unmarshal(data []byte, v interface{}) error {
	// 设置 useNumber = true  这样 json.Unmarshal 能够正常处理长整型
	var k = jsoniter.Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true, // will lose precession
		ObjectFieldMustBeSimpleString: true, // do not unescape object field
		UseNumber:                     true,
	}.Froze()
	return k.Unmarshal(data, v)
}
