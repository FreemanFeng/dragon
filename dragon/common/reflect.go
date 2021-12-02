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
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func TypeOf(v interface{}) reflect.Kind {
	return reflect.TypeOf(v).Kind()
}

func CheckType(v interface{}, expect reflect.Kind) bool {
	if v == nil {
		return false
	}
	if TypeOf(v) == expect {
		return true
	}
	return false
}

func IsMap(v interface{}) bool {
	return CheckType(v, reflect.Map)
}

func IsMapBytes(msg interface{}) bool {
	if !IsMap(msg) {
		return false
	}
	if _, ok := msg.(map[string][]byte); !ok {
		return false
	}
	return true
}

func CheckTypes(v interface{}, expects []reflect.Kind) bool {
	for _, k := range expects {
		if CheckType(v, k) {
			return true
		}
	}
	return false
}

func IsList(v interface{}) bool {
	types := []reflect.Kind{reflect.Array, reflect.Slice}
	return CheckTypes(v, types)
}

func IsListBytes(msg interface{}) bool {
	if !IsList(msg) {
		return false
	}
	if _, ok := msg.([][]byte); !ok {
		return false
	}
	return true
}

func IsInt(v interface{}) bool {
	types := []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}
	return CheckTypes(v, types)
}

func IsUInt(v interface{}) bool {
	types := []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}
	return CheckTypes(v, types)
}

func IsFloat(v interface{}) bool {
	types := []reflect.Kind{reflect.Float32, reflect.Float64}
	return CheckTypes(v, types)
}

func IsBool(v interface{}) bool {
	return CheckType(v, reflect.Bool)
}

func IsNull(v interface{}) bool {
	return v == nil
}

func IsString(v interface{}) bool {
	return CheckType(v, reflect.String)
}

func IsBytes(v interface{}) bool {
	if _, ok := v.([]byte); ok {
		return true
	}
	return false
}

func IfToString(v interface{}) string {
	if v == nil {
		return EMPTY
	}
	switch TypeOf(v) {
	case reflect.Int:
		return strconv.Itoa(v.(int))
	case reflect.Int8:
		return strconv.Itoa(int(v.(int8)))
	case reflect.Int16:
		return strconv.Itoa(int(v.(int16)))
	case reflect.Int32:
		return strconv.Itoa(int(v.(int32)))
	case reflect.Int64:
		return strconv.FormatInt(v.(int64), 10)
	case reflect.Uint8:
		return strconv.Itoa(int(v.(uint8)))
	case reflect.Uint16:
		return strconv.Itoa(int(v.(uint16)))
	case reflect.Uint32:
		return strconv.FormatUint(uint64(v.(uint32)), 10)
	case reflect.Uint64:
		return strconv.FormatUint(uint64(v.(uint64)), 10)
	case reflect.Float32:
		return fmt.Sprintf("%f", v.(float32))
	case reflect.Float64:
		return fmt.Sprintf("%f", v.(float64))
	case reflect.Bool:
		switch v.(bool) {
		case true:
			return "true"
		case false:
			return "false"
		}
	case reflect.String:
		return ToString(v)
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return ToString(v)
}

func IfToBytes(v interface{}) []byte {
	if v == nil {
		return []byte(EMPTY)
	}
	if IsList(v) || IsMap(v) {
		b, err := json.Marshal(v)
		if err != nil {
			Err(ModuleCommon, "!!!!!!序列化", v, "出错", err)
			return ToBytes(v)
		}
		return b
	}
	switch TypeOf(v) {
	case reflect.Int:
		return []byte(strconv.Itoa(v.(int)))
	case reflect.Int8:
		return []byte(strconv.Itoa(int(v.(int8))))
	case reflect.Int16:
		return []byte(strconv.Itoa(int(v.(int16))))
	case reflect.Int32:
		return []byte(strconv.Itoa(int(v.(int32))))
	case reflect.Int64:
		return []byte(strconv.FormatInt(v.(int64), 10))
	case reflect.Uint8:
		return []byte(strconv.Itoa(int(v.(uint8))))
	case reflect.Uint16:
		return []byte(strconv.Itoa(int(v.(uint16))))
	case reflect.Uint32:
		return []byte(strconv.FormatUint(uint64(v.(uint32)), 10))
	case reflect.Uint64:
		return []byte(strconv.FormatUint(uint64(v.(uint64)), 10))
	case reflect.Float32:
		return []byte(fmt.Sprintf("%f", v.(float32)))
	case reflect.Float64:
		return []byte(fmt.Sprintf("%f", v.(float64)))
	case reflect.Bool:
		switch v.(bool) {
		case true:
			return []byte("true")
		case false:
			return []byte("false")
		}
	case reflect.String:
		return []byte(v.(string))
	}
	if b, ok := v.([]byte); ok {
		return b
	}
	return ToBytes(v)
}

func IfToBase64(v interface{}) string {
	b := IfToBytes(v)
	return Base64Encode(b)
}

func IfToInt(v interface{}) int64 {
	if v == nil || IsList(v) || IsMap(v) {
		return 0
	}
	switch TypeOf(v) {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Int64:
		return v.(int64)
	case reflect.Uint8:
		return int64(v.(uint8))
	case reflect.Uint16:
		return int64(v.(uint16))
	case reflect.Uint32:
		return int64(v.(uint32))
	case reflect.Uint64:
		return int64(v.(uint64))
	case reflect.Float32:
		x := v.(float32)
		k := math.Floor(float64(x))
		return int64(k)
	case reflect.Float64:
		x := v.(float64)
		k := math.Floor(x)
		return int64(k)
	case reflect.Bool:
		switch v.(bool) {
		case true:
			return 1
		case false:
			return 0
		}
	case reflect.String:
		return 0
	}
	if _, ok := v.([]byte); ok {
		return 0
	}
	return 0
}

func IfToUInt(v interface{}) uint64 {
	if v == nil || IsList(v) || IsMap(v) {
		return 0
	}
	switch TypeOf(v) {
	case reflect.Int:
		return uint64(v.(int))
	case reflect.Int8:
		return uint64(v.(int8))
	case reflect.Int16:
		return uint64(v.(int16))
	case reflect.Int32:
		return uint64(v.(int32))
	case reflect.Int64:
		return uint64(v.(int64))
	case reflect.Uint8:
		return uint64(v.(uint8))
	case reflect.Uint16:
		return uint64(v.(uint16))
	case reflect.Uint32:
		return uint64(v.(uint32))
	case reflect.Uint64:
		return uint64(v.(uint64))
	case reflect.Float32:
		x := v.(float32)
		k := math.Floor(float64(x))
		return uint64(k)
	case reflect.Float64:
		x := v.(float64)
		k := math.Floor(x)
		return uint64(k)
	case reflect.Bool:
		switch v.(bool) {
		case true:
			return 1
		case false:
			return 0
		}
	case reflect.String:
		return 0
	}
	if _, ok := v.([]byte); ok {
		return 0
	}
	return 0
}

func IfToFloat(v interface{}) float64 {
	if v == nil || IsList(v) || IsMap(v) {
		return 0.0
	}
	switch TypeOf(v) {
	case reflect.Int:
		return float64(v.(int))
	case reflect.Int8:
		return float64(v.(int8))
	case reflect.Int16:
		return float64(v.(int16))
	case reflect.Int32:
		return float64(v.(int32))
	case reflect.Int64:
		return float64(v.(int64))
	case reflect.Uint8:
		return float64(v.(uint8))
	case reflect.Uint16:
		return float64(v.(uint16))
	case reflect.Uint32:
		return float64(v.(uint32))
	case reflect.Uint64:
		return float64(v.(uint64))
	case reflect.Float32:
		return float64(v.(float32))
	case reflect.Float64:
		return v.(float64)
	case reflect.Bool:
		switch v.(bool) {
		case true:
			return 1.0
		case false:
			return 0.0
		}
	case reflect.String:
		return 0.0
	}
	if _, ok := v.([]byte); ok {
		return 0.0
	}
	return 0.0
}

func IfToJsonByte(v interface{}) [][]byte {
	var bs [][]byte
	if IsString(v) {
		return append(bs, []byte("STR"), IfToBytes(v))
	}
	if IsBytes(v) {
		return append(bs, []byte("BYTES"), v.([]byte))
	}
	if IsInt(v) {
		return append(bs, []byte("INT"), IfToBytes(v))
	}
	if IsUInt(v) {
		return append(bs, []byte("UINT"), IfToBytes(v))
	}
	if IsFloat(v) {
		return append(bs, []byte("FLOAT"), IfToBytes(v))
	}
	if IsBool(v) {
		return append(bs, []byte("BOOL"), IfToBytes(v))
	}
	if IsList(v) {
		b, e := json.Marshal(v)
		if e == nil {
			return append(bs, []byte("LIST"), b)
		}
	}

	if IsMap(v) {
		b, e := json.Marshal(v)
		if e == nil {
			return append(bs, []byte("MAP"), b)
		}
	}
	return append(bs, []byte("NULL"), nil)
}

func IfToJsonStr(v interface{}) []string {
	var s []string
	if IsString(v) {
		return append(s, STR, IfToString(v))
	}
	if IsBytes(v) {
		return append(s, BYTES, Base64Encode(v.([]byte)))
	}
	if IsInt(v) {
		return append(s, INT, IfToString(v))
	}
	if IsUInt(v) {
		return append(s, UINT, IfToString(v))
	}
	if IsFloat(v) {
		return append(s, FLOAT, IfToString(v))
	}
	if IsBool(v) {
		return append(s, BOOL, IfToString(v))
	}

	if IsListBytes(v) {
		var r []string
		for _, b := range v.([][]byte) {
			r = append(r, Base64Encode(b))
		}
		b, e := json.Marshal(r)
		if e == nil {
			return append(s, ListBytes, Base64Encode(b))
		}
	}

	if IsList(v) {
		b, e := json.Marshal(v)
		if e == nil {
			return append(s, LIST, Base64Encode(b))
		}
	}

	if IsMapBytes(v) {
		r := map[string]string{}
		for k, b := range v.(map[string][]byte) {
			r[k] = Base64Encode(b)
		}
		b, e := json.Marshal(r)
		if e == nil {
			return append(s, MapBytes, Base64Encode(b))
		}
	}

	if IsMap(v) {
		b, e := json.Marshal(v)
		if e == nil {
			return append(s, MAP, Base64Encode(b))
		}
	}

	return append(s, NULL, NULL)
}

func AlignType(field, value interface{}) interface{} {
	if field == nil || value == nil {
		return nil
	}
	if TypeOf(field) == TypeOf(value) {
		return value
	}
	v := IfToString(value)
	if IsString(field) {
		if strings.HasPrefix(v, DoubleQuotes) || strings.HasSuffix(v, DoubleQuotes) {
			return strings.TrimSuffix(strings.TrimPrefix(v, DoubleQuotes), DoubleQuotes)
		}
		return v
	}
	if IsInt(field) {
		k, _ := strconv.ParseInt(v, 10, 64)
		return k
	}
	if IsUInt(field) {
		k, _ := strconv.ParseUint(v, 10, 64)
		return k
	}
	if IsFloat(field) {
		k, _ := strconv.ParseFloat(v, 64)
		return k
	}
	if IsBool(field) {
		if v == "true" {
			return true
		}
		return false
	}
	return nil
}

func DecodeJsonList(b []byte, p *PluginType) []map[string][]byte {
	if IsGoCode(p) {
		return DecodeJsonListByte(b)
	}
	return DecodeJsonListStr(b)
}

func DecodeJsonListStr(b []byte) []map[string][]byte {
	var msg []string
	var ms []map[string][]byte
	e := json.Unmarshal(b, &msg)
	if e != nil {
		Err(ModuleProto, "1.反序列化出错", e)
		return nil
	}
	for _, data := range msg {
		//Info(ModuleProto, "JSON列表元素类型", TypeOf(data))
		var m map[string]interface{}
		r := map[string][]byte{}
		x := Base64Decode(data)
		e = json.Unmarshal(x, &m)
		if e != nil {
			Err(ModuleProto, "2. 反序列化出错", e)
			return nil
		}
		//Info(ModuleProto, "m:", m)
		for k, v := range m {
			//Info(ModuleProto, "K类型", TypeOf(k), "V类型", TypeOf(v))
			b, e := json.Marshal(v)
			if e != nil {
				Err(ModuleProto, "3. 反序列化出错", e)
				return nil
			}
			r[k] = b
		}
		ms = append(ms, r)
	}
	//Info(ModuleProto, "ms:", ms)
	return ms
}

func DecodeJsonListByte(b []byte) []map[string][]byte {
	var bs [][]byte
	var ms []map[string][]byte
	e := json.Unmarshal(b, &bs)
	if e != nil {
		Err(ModuleProto, "反序列化出错", e)
		return nil
	}
	for _, v := range bs {
		var m map[string][]byte
		e = json.Unmarshal(v, &m)
		if e != nil {
			Err(ModuleProto, "反序列化出错", e)
			return nil
		}
		ms = append(ms, m)
	}
	return ms
}
