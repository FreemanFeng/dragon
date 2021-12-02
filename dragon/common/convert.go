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
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func JoinKeys(sep string, keys ...interface{}) string {
	n := len(keys)
	s := make([]string, n)
	for i, k := range keys {
		s[i] = IfToString(k)
	}
	return strings.Join(s, sep)
}

func Join(sep string, keys ...string) string {
	return strings.Join(keys, sep)
}

func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func ToBytes(v interface{}) []byte {
	return []byte(ToString(v))
}

func SortMap(x interface{}) string {
	var r []string
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Map {
		return ToString(x)
	}
	for _, k := range v.MapKeys() {
		t := v.MapIndex(k)
		r = append(r, JoinKeys(COLON, k.Interface(), t.Interface()))
	}
	sort.Strings(r)
	return strings.Join(r, SPACE)
}

func AddBackSlash(text, char string) string {
	s := []string{BACKSLASH, char}
	h := strings.Join(s, EMPTY)
	return strings.Replace(text, char, h, UNDEFINED)
}

func NormalizeRegexp(text string) string {
	t := text
	s := []string{"\\", ".", "+", "*", "?", "(", ")", "|", "[", "]", "{", "}", "^", "$"}
	for _, k := range s {
		t = AddBackSlash(t, k)
	}
	return t
}

func TrimSpaces(c []byte) []byte {
	var bs [][]byte
	b := ReplaceNoChars(c)
	b = TrimSpacesAroundLogicPlus(b)
	b = TrimSpacesAroundLogicSub(b)
	b = TrimSpacesAroundLogicMulti(b)
	b = TrimSpacesAroundLP(b)
	b = TrimSpacesAroundRP(b)
	b = TrimSpacesAroundMod(b)
	b = TrimSpacesAroundPower(b)
	b = TrimSpacesAroundLS(b)
	b = TrimSpacesAroundRS(b)
	b = TrimSpacesAroundLogicOr(b)
	b = TrimSpacesAroundLogicAnd(b)
	b = TrimSpacesAroundLogicNot(b)
	b = TrimSpacesAroundLogicEOR(b)
	b = TrimSpacesAroundLogicDiv(b)
	b = TrimSpacesAroundColon(b)
	b = TrimSpacesAroundDot(b)
	b = TrimSpacesAroundDotDot(b)
	b = TrimSpacesAroundDotDotDot(b)
	b = ConvertSpacesInsideQuotes(b)
	b = TrimSpacesAroundEqual(b)
	b = TrimSpacesAroundNotEqual(b)
	b = TrimSpacesAroundMoreThan(b)
	b = TrimSpacesAroundLessThan(b)
	b = TrimSpacesAroundSetVar(b)
	b = TrimSpacesAroundNotPresent(b)
	b = TrimSpacesAroundExtraPresent(b)
	b = TrimSpacesAroundOr(b)
	b = TrimSpacesAroundOn(b)
	b = TrimSpacesAroundWithout(b)
	b = TrimSpacesAroundWith(b)
	rs := regexp.MustCompile(`\s+`)
	b = rs.ReplaceAll(bytes.TrimSpace(b), []byte(SPACE))
	k := bytes.Split(b, []byte(SPACE))
	last := []byte(EMPTY)
	cops := CombineOps()
	bops := BasicOps()
	for i, v := range k {
		if i == 0 {
			bs = append(bs, v)
			last = v
			continue
		}
		sep := []byte(SPACE)
		if !MatchInside(v, bops...) && !MatchInside(v, cops...) && !MatchInside(last, cops...) {
			sep = []byte(COMMA)
		} else if MatchInside(v, MoreThan) && MatchInside(v, OpRS) {
			sep = []byte(COMMA)
		} else if MatchInside(v, LessThan) && MatchInside(v, OpLS) {
			sep = []byte(COMMA)
		} else {
			last = v
		}
		bs = append(bs, sep, v)
	}
	b = bytes.Join(bs, []byte(EMPTY))
	//con := []byte(EMPTY)
	// 在构造操作里(->)，通过:定义模板变量，此时:不作为并发操作符
	if MatchInside(b, OpConstruct) && MatchInside(b, OpColon) {
		//Log(ModuleParser, "构造变量", string(b))
		i := bytes.Index(b, []byte(OpColon))
		k := bytes.Split(b[:i], []byte(DOT))
		for _, v := range k {
			if len(v) == 0 {
				continue
			}
			_, err := strconv.Atoi(string(v))
			if err != nil {
				//Log("转换出错", v, err)
				i = 0
				break
			}
		}
		if i > 0 {
			//con = TrimSpacesAroundOp(b[:i+1], OpColon, OpColon)
			//Log("并发前缀", string(con))
			b = b[i+1:]
		}
		b = TrimSpacesAroundOp(b, OpColon, OpColon)
	}
	//Log(ModuleParser, "拼接前", string(b))
	// TODO not support yet
	//if len(con) > 0 {
	//	b = bytes.Join([][]byte{con, b}, []byte(SPACE))
	//}
	//Log(ModuleParser, "拼接后", string(b))
	if MatchEnd(b, RightBracket) {
		n := bytes.LastIndex(b, []byte(LeftBracket))
		if n > 1 {
			if bytes.Equal(b[n-1:n], []byte(OpEqual)) || bytes.Equal(b[n-1:n], []byte(OpConstruct)) {
				k := len(b)
				t := bytes.Trim(b[n+1:k-1], COMMA)
				bs = [][]byte{b[:n+1], []byte(COMMA), t, []byte(COMMA), []byte(RightBracket)}
				b = bytes.Join(bs, []byte(EMPTY))
				return b
			}
		}
	}
	return b
}

func ConvertSpacesInsideQuotes(c []byte) []byte {
	k := bytes.Index(c, []byte(DoubleQuotes))
	if k == UNDEFINED {
		return c
	}
	n := len(c)
	var r []byte
	mark := byte(SpaceChar)
	b := []byte(BoundaryChars)
	for i := 0; i < n; i++ {
		switch {
		case c[i] == DoubleQuotesChar || c[i] == SingleQuotesChar:
			r = append(r, c[i])
			switch mark {
			case c[i]:
				mark = byte(SpaceChar)
			case byte(SpaceChar):
				mark = c[i]
			default:
				r = append(r, c[i])
			}
		case c[i] == byte(SpaceChar) && mark != byte(SpaceChar):
			for _, v := range b {
				r = append(r, v)
			}
		default:
			r = append(r, c[i])
		}
	}
	return r
}

func TrimSpacesAroundEqual(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpEqual, OpEqual)
}

func TrimSpacesAroundNotEqual(c []byte) []byte {
	return TrimSpacesAroundOp(c, NotEqual, NotEqual)
}

func TrimSpacesAroundMoreThan(c []byte) []byte {
	return TrimSpacesAroundOp(c, MoreThan, MoreThan)
}

func TrimSpacesAroundLessThan(c []byte) []byte {
	return TrimSpacesAroundOp(c, LessThan, LessThan)
}

func TrimSpacesAroundSetVar(c []byte) []byte {
	return TrimSpacesAroundExtraOp(c, OpConstruct)
}

func TrimSpacesAroundNotPresent(c []byte) []byte {
	return TrimSpacesAroundExtraOp(c, NotPresent)
}

func TrimSpacesAroundExtraPresent(c []byte) []byte {
	return TrimSpacesAroundExtraOp(c, ExtraPresent)
}

func ReplaceNoChars(c []byte) []byte {
	return bytes.Replace(c, []byte(DoubleQuotes2), []byte(NoChars), UNDEFINED)
}

func TrimSpacesAroundLogicPlus(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpPlus, OpPlus)
}

func TrimSpacesAroundLogicSub(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpSub, OpSub)
}

func TrimSpacesAroundLogicMulti(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpMulti, OpMulti)
}

func TrimSpacesAroundLP(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpLP, OpLP)
}

func TrimSpacesAroundRP(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpRP, OpRP)
}

func TrimSpacesAroundMod(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpMod, OpMod)
}

func TrimSpacesAroundPower(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpPower, OpPower)
}
func TrimSpacesAroundLS(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpLS, OpLS)
}
func TrimSpacesAroundRS(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpRS, OpRS)
}
func TrimSpacesAroundLogicOr(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpOr, OpOr)
}
func TrimSpacesAroundLogicAnd(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpAnd, OpAnd)
}

func TrimSpacesAroundLogicNot(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpNot, OpNot)
}
func TrimSpacesAroundLogicEOR(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpEOR, OpEOR)
}

func TrimSpacesAroundLogicDiv(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpDiv, OpDiv)
}

func TrimSpacesAroundColon(c []byte) []byte {
	return TrimSpacesAroundExtraOp(c, OpColon)
}

func TrimSpacesAroundDot(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpDot, OpDot)
}

func TrimSpacesAroundDotDot(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpDotDot, OpDotDot)
}

func TrimSpacesAroundDotDotDot(c []byte) []byte {
	return TrimSpacesAroundOp(c, OpDotDotDot, OpDotDotDot)
}

func TrimSpacesAroundOn(c []byte) []byte {
	return TrimSpacesAroundTextOp(c, OpOn)
}

func TrimSpacesAroundOr(c []byte) []byte {
	return TrimSpacesAroundTextOp(c, OR)
}

func TrimSpacesAroundWithout(c []byte) []byte {
	return TrimSpacesAroundTextOp(c, OpWithout)
}

func TrimSpacesAroundWith(c []byte) []byte {
	return TrimSpacesAroundTextOp(c, OpWith)
}

func TrimSpacesWithRegexp(c []byte, rs, opb string) []byte {
	rq := regexp.MustCompile(rs)
	if rq.Match(c) {
		k := rq.ReplaceAll(bytes.TrimSpace(c), []byte(opb))
		return k
	}
	return c
}

func TrimSpacesAroundOp(c []byte, opa, opb string) []byte {
	s := []string{`\s*`, `\s*`}
	h := strings.Join(s, NormalizeRegexp(opa))
	return TrimSpacesWithRegexp(c, h, opb)
}

func TrimSpacesAroundExtraOp(c []byte, op string) []byte {
	s := []string{EMPTY, SPACE}
	h := strings.Join(s, op)
	return TrimSpacesAroundOp(c, op, h)
}

func TrimSpacesAroundTextOp(c []byte, op string) []byte {
	s := []string{`\s+`, `\s+`}
	k := []string{SPACE, SPACE}
	h := strings.Join(s, NormalizeRegexp(op))
	t := strings.Join(k, op)
	return TrimSpacesWithRegexp(c, h, t)
}

func SplitOps(c []byte, ot []OpType) []OpType {
	n := len(ot)
	ot = SplitColon(c, ot)
	ot = SplitSetVar(c, ot)
	ot = SplitNotPresent(c, ot)
	ot = SplitExtraPresent(c, ot)
	ot = SplitEqualOrNotEqual(c, ot)
	ot = SplitMoreThan(c, ot)
	ot = SplitLessThan(c, ot)
	ot = SplitWithout(c, ot)
	ot = SplitWith(c, ot)
	ot = SplitOr(c, ot)
	ot = SplitOn(c, ot)
	if len(ot) == n && len(bytes.TrimSpace(c)) > 0 {
		k := OpType{}
		splitRange(&k, bytes.TrimSpace(c))
		if n > 0 {
			ot[n-1].Value = append(ot[n-1].Value, k.Value...)
		} else {
			ot = append(ot, k)
		}
	}
	return ot
}

func SplitEqualOrNotEqual(c []byte, ot []OpType) []OpType {
	if MatchInside(c, NotEqual) {
		return SplitOp(c, NotEqual, ot)
	}
	return SplitOp(c, OpEqual, ot)
}

func SplitMoreThan(c []byte, ot []OpType) []OpType {
	if MatchInside(c, OpRS, OpConstruct) {
		return ot
	}
	return SplitOp(c, MoreThan, ot)
}

func SplitLessThan(c []byte, ot []OpType) []OpType {
	return SplitOp(c, LessThan, ot)
}

func SplitSetVar(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OpConstruct, ot)
}

func SplitColon(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OpColon, ot)
}

func SplitNotPresent(c []byte, ot []OpType) []OpType {
	return SplitOp(c, NotPresent, ot)
}

func SplitExtraPresent(c []byte, ot []OpType) []OpType {
	return SplitOp(c, ExtraPresent, ot)
}

func SplitOr(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OR, ot)
}

func SplitOn(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OpOn, ot)
}

func SplitWithout(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OpWithout, ot)
}

func SplitWith(c []byte, ot []OpType) []OpType {
	return SplitOp(c, OpWith, ot)
}

func splitKey(k *OpType, b []byte) {
	t := bytes.Split(b, []byte(OpDot))
	for _, v := range t {
		if string(v) == EMPTY {
			continue
		}
		_, e := strconv.Atoi(string(v))
		if e != nil {
			return
		}
	}
	d := 0
	n := 0
	last := 0
	for _, v := range t {
		if string(v) == EMPTY {
			d = UNDEFINED
			continue
		}
		n, _ = strconv.Atoi(string(v))
		if d == UNDEFINED {
			for i := last + 1; i <= n; i++ {
				k.CIDs = append(k.CIDs, i)
			}
			d = 0
		} else {
			k.CIDs = append(k.CIDs, n)
		}
		last = n
	}
}

func splitRange(k *OpType, b []byte) {
	p := bytes.Split(b, []byte(COMMA))
	for _, s := range p {
		x := bytes.Split(s, []byte(OpDotDot))
		if len(x) > 1 {
			if MatchInside(s, OpDotDotDot, OpRP, OpLP, OpMulti, OpPlus, OpSub, OpDiv,
				OpMod, OpPower, OpLS, OpRS, OpOr, OpAnd, OpEOR, OpNot) {
				k.Value = append(k.Value, string(s))
				continue
			}
			Debug(ModuleCommon, "未处理范围", string(s))
			v1, e1 := strconv.Atoi(string(x[0]))
			v2, e2 := strconv.Atoi(string(x[1]))
			if e1 == nil && e2 == nil && v1 <= v2 {
				Debug(ModuleCommon, "范围数据", v1, v2)
				for i := v1; i <= v2; i++ {
					k.Value = append(k.Value, strconv.Itoa(i))
				}
				continue
			}
		}
		k.Value = append(k.Value, string(s))
	}
}

func SplitOp(c []byte, opa string, ot []OpType) []OpType {
	// 组合操作符
	if string(c) == opa {
		k := OpType{Op: opa, Continue: true}
		n := len(ot)
		for n > 0 {
			ot[n-1].Continue = true
			n -= 1
		}
		return append(ot, k)
	} else if MatchInside(c, opa) && !Equal(opa, OpOn, OpWith, OpWithout, OR) {
		Debug(ModuleCommon, "操作符", opa, "待处理内容", string(c))
		t := bytes.Split(bytes.TrimSpace(c), []byte(opa))
		k := OpType{Key: string(t[0]), Op: opa}
		splitKey(&k, bytes.TrimSpace(t[0]))
		// 模板变量，非并发设置
		if opa == OpColon && len(k.CIDs) == 0 {
			return ot
		}
		// 并发设置，key为空
		if opa == OpColon {
			k.Key = EMPTY
		}
		// 组合操作符
		cops := CombineOps()
		n := len(ot)
		if Equal(opa, cops...) {
			k.Continue = true
			for n > 0 {
				ot[n-1].Continue = true
				n -= 1
			}
			return append(ot, k)
		}
		v := bytes.TrimSpace(t[1])
		Debug(ModuleCommon, "未处理内容", string(v), "操作符", opa)
		splitRange(&k, v)
		if n > 0 && ot[n-1].Continue {
			k.Continue = true
		}
		return append(ot, k)
	}
	return ot
}

func BuildFormat(s string) FormatType {
	f := FormatType{Content: s}
	_, e := strconv.Atoi(s)
	if strings.Index(s, DoubleQuotes) == 0 {
		f.MustString = true
		f.Content = TrimDoubleQuotes(s)
	} else if e == nil {
		f.IsInt = true
	}
	return f
}

func TrimDoubleQuotes(s string) string {
	if strings.Index(s, DoubleQuotes) == 0 {
		return strings.TrimSuffix(strings.TrimPrefix(s, DoubleQuotes), DoubleQuotes)
	}
	return s
}

func TrimOp(op *OpType) {
	op.Key = TrimDoubleQuotes(op.Key)
	for i := range op.Value {
		op.Value[i] = TrimDoubleQuotes(op.Value[i])
	}
}

func ConvertValue(v string) interface{} {
	if strings.HasPrefix(v, DoubleQuotes) || strings.HasSuffix(v, DoubleQuotes) {
		return strings.TrimSuffix(strings.TrimPrefix(v, DoubleQuotes), DoubleQuotes)
	}
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	k, e := strconv.Atoi(v)
	if e == nil {
		return k
	}
	i, e := strconv.ParseInt(v, 10, 64)
	if e == nil {
		return i
	}
	d, e := strconv.ParseFloat(v, 64)
	if e == nil {
		return d
	}
	return v
}

func JsonListByte(params ...interface{}) []byte {
	var bs [][]byte
	for _, v := range params {
		bs = append(bs, IfToJsonByte(v)...)
	}
	b, e := json.Marshal(bs)
	if e != nil {
		Err(ModuleCommon, "封装json列表出错", e)
		return nil
	}
	return b
}
func JsonListStr(params ...interface{}) []byte {
	var s []string
	for _, v := range params {
		s = append(s, IfToJsonStr(v)...)
	}
	b, e := json.Marshal(s)
	if e != nil {
		Err(ModuleCommon, "封装json列表出错", e)
		return nil
	}
	return b
}

func LowerCaseNoHyphen(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, HYPHEN, UnderScope))
}

func AddUnderScope(s string) string {
	k := []string{UnderScope, s}
	return strings.Join(k, EMPTY)
}

func SplitKeys(s string, keys ...string) []string {
	for _, k := range keys {
		if strings.Contains(s, k) {
			return strings.Split(s, k)
		}
	}
	return []string{s}
}

func CopyVarData(key string, src *VarData, dst map[string]*VarData) {
	if dst[key] == nil {
		dst[key] = &VarData{ID: src.ID, Title: src.Title, Host: src.Host, Total: src.Total,
			Data: map[string][]*FileType{}, Tags: map[string]*OpType{}, Funcs: map[string]*PluginType{},
			Tops: map[string][]string{}, Msgs: map[string]interface{}{},
			Keys: map[string][]FormatType{}, Bins: map[string]string{}, Vars: map[string]interface{}{},
			VPos: map[string]int{}, Ctx: map[string]interface{}{}}
	}
	dp := dst[key]
	for k, v := range src.Data {
		for _, item := range v {
			p := &FileType{}
			*p = *item
			copy(p.Tags, item.Tags)
			copy(p.Content, item.Content)
			dp.Data[k] = append(dp.Data[k], p)
		}
	}
	for k, v := range src.Tags {
		p := &OpType{}
		*p = *v
		copy(p.CIDs, v.CIDs)
		copy(p.Value, v.Value)
		dp.Tags[k] = p
	}
	// 共享
	for k, v := range src.Funcs {
		dp.Funcs[k] = v
	}
	for k, v := range src.Tops {
		copy(dp.Tops[k], v)
	}
	for k, v := range src.Msgs {
		dp.Msgs[k] = v
	}
	for k, v := range src.Keys {
		copy(dp.Keys[k], v)
	}
	for k, v := range src.Bins {
		dp.Bins[k] = v
	}
	for k, v := range src.Vars {
		dp.Vars[k] = v
	}
	for k, v := range src.VPos {
		dp.VPos[k] = v
	}
	for k, v := range src.Ctx {
		dp.Ctx[k] = v
	}
}

func CopyVarDataMap(src, dst map[string]*VarData) {
	for k, v := range src {
		dst[k] = nil
		CopyVarData(k, v, dst)
	}
}

func SplitVars(s string) []string {
	var rs []string
	start := 0
	i := 0
	n := len(s)
	//Log("n:", n, "s:", s)
	for i < n {
		//Log("i:", i, "start:", start)
		if StrMatchStart(s[i:i+1], LeftParentheses, RightParentheses, OpPlus, OpSub, OpMulti, OpDiv, OpMod) {
			if i == start {
				//Log(">>>>", s[i:i+1])
				rs = append(rs, s[i:i+1])
			} else {
				//Log(">>>>", s[start:i])
				rs = append(rs, s[start:i])
				//Log(">>>>", s[i:i+1])
				rs = append(rs, s[i:i+1])
			}
			start = i + 1
		}
		i++
	}
	if start < n {
		//Log(">>>>", s[start:n])
		rs = append(rs, s[start:n])
	}
	return rs
}

func IsOpVars(s string) bool {
	return StrMatch(s, LeftParentheses, RightParentheses, OpPlus, OpSub, OpMulti, OpDiv, OpMod)
}

func ParseVar(s string, vd *VarData) string {
	if IsOpVars(s) {
		return s
	}
	k, ok := vd.Vars[s]
	if ok {
		return IfToString(k)
	}
	return s
}

func EvalVarData(key, s string, vd *VarData) string {
	if !IsOpVars(s) {
		return s
	}
	var ks []string
	for _, k := range SplitVars(s) {
		ks = append(ks, ParseVar(k, vd))
	}
	h := strings.Join(ks, EMPTY)
	r := Calculate(h)
	if r == h {
		return s
	}
	vd.Vars[key] = r
	return r
}

func IsDigit(v string) bool {
	_, err := strconv.Atoi(v)
	if err == nil {
		return true
	}
	_, err = strconv.ParseInt(v, 10, 64)
	if err == nil {
		return true
	}
	_, err = strconv.ParseFloat(v, 64)
	if err == nil {
		return true
	}
	return false
}

func LoadKeys(cd *VarData, vt map[string]*VarType, s string) {
	if s == EMPTY {
		return
	}
	keys := BuildKeys(s)
	key := keys[0].Content
	_, ok := vt[key]
	if !ok {
		return
	}
	_, ok = cd.Keys[s]
	if ok {
		return
	}
	Info(ModuleCommon, ">>>>>>> key", s, "已映射为", keys)
	cd.Keys[s] = keys
}
