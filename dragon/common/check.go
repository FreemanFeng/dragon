package common

import (
	"sort"
	"strconv"
	"strings"
)

func CompareDigits(expect, op, receive string) bool {
	k1, e1 := strconv.ParseFloat(receive, 64)
	k2, e2 := strconv.ParseFloat(expect, 64)
	if e1 != nil || e2 != nil {
		Info(ModuleCommon, "校验", receive, op, expect, "失败!")
		return false
	}
	if k1 > k2 && op == MoreThan {
		Info(ModuleCommon, "校验", receive, op, expect, "成功!")
		return true
	} else if k1 < k2 && op == LessThan {
		Info(ModuleCommon, "校验", receive, op, expect, "成功!")
		return true
	}
	return false
}

func CheckExpected(expect interface{}, op string, receive interface{}) bool {
	if receive == nil {
		Err(ModuleCommon, receive, "为nil，校验", receive, op, expect, "失败!")
		return false
	}
	if k, ok := expect.(string); ok {
		expect = k
	}
	if expect == strings.TrimSpace(NoChars) {
		Info(ModuleCommon, ">>>>>>>>>>>> 校验前替换", expect, "为空")
		expect = EMPTY
	}
	Info(ModuleCommon, "准备校验", receive, op, expect, "是否成功?")
	// 字典校验
	if k, ok := receive.(map[string]interface{}); ok {
		if Equal(op, OpWith, OpWithout) {
			r := true
			if op == OpWithout {
				r = false
			}
			_, ok = k[IfToString(expect)]
			Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "结果为", ok == r)
			return ok == r
		}
		if x, ok := expect.(map[string]interface{}); ok {
			var ks, xs []string
			for _, v := range k {
				ks = append(ks, IfToString(v))
			}
			for _, v := range x {
				xs = append(xs, IfToString(v))
			}
			sort.Strings(ks)
			sort.Strings(xs)
			if len(ks) != len(xs) || IfToString(ks) != IfToString(xs) {
				Err(ModuleCommon, "收到的值列表", ks, "不等于期望的值列表", xs)
				return false
			}
			Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "成功")
			return true
		} else {
			Err(ModuleCommon, "校验", receive, op, expect, "失败")
			return false
		}
	}
	// 列表校验
	if k, ok := receive.([]interface{}); ok {
		if Equal(op, OpWith, OpWithout) {
			found := false
			for _, v := range k {
				if IfToString(v) == IfToString(expect) {
					if op == OpWith {
						Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "成功")
						return true
					}
					found = true
				}
			}
			if op == OpWithout {
				if !found {
					Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "成功")
					return true
				}
				Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "失败")
				return false
			}
		}
		if x, ok := expect.([]interface{}); ok {
			var ks, xs []string
			for _, v := range k {
				ks = append(ks, IfToString(v))
			}
			for _, v := range x {
				xs = append(xs, IfToString(v))
			}
			sort.Strings(ks)
			sort.Strings(xs)
			if len(ks) != len(xs) || IfToString(ks) != IfToString(xs) {
				Err(ModuleCommon, "收到的值列表", ks, "不等于期望的值列表", xs)
				return false
			}
			Info(ModuleCommon, ">>>>>>> 校验", receive, op, expect, "成功")
			return true
		} else {
			Err(ModuleCommon, "校验", receive, op, expect, "失败")
			return false
		}
	}
	exp := IfToString(expect)
	rev := IfToString(receive)
	if op == MoreThan || op == LessThan {
		return CompareDigits(exp, op, rev)
	}
	// 字符串相等 或 数字相等
	if op == OpEqual && expect == rev {
		Info(ModuleCommon, "校验", receive, op, expect, "成功!")
		return true
	}
	// 接收到的字符串是否包含期望的字符串
	if op == OpWith && strings.Index(rev, exp) >= 0 {
		Info(ModuleCommon, "校验", receive, op, expect, "成功!")
		return true
	}
	// 接收到的字符串是否不包含期望的字符串
	if op == OpWithout && strings.Index(rev, exp) == UNDEFINED {
		Info(ModuleCommon, "校验", receive, op, expect, "成功!")
		return true
	}
	Err(ModuleCommon, "校验", receive, op, expect, "失败!")
	return false
}
