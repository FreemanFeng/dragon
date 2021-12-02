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

package builtin

import (
	"math/rand"
	"time"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func Run() {
	rand.Seed(time.Now().Unix())
	//日期时间
	SetCall(Join(COLON, ModuleBuiltin, FuncNameFormatDate), FormatDate)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameDATE), DATE)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameFormatTime), FormatTime)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameTIME), TIME)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameFormatSeconds), FormatSeconds)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameSEC), SEC)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameFormatNowInSeconds), FormatNowInSeconds)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameNOW), NOW)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameNowAddSeconds), NowAddSeconds)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameNowMS), NowMS)
	//编码
	SetCall(Join(COLON, ModuleBuiltin, FuncNameMD5), MD5)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameSHA1), SHA1)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameBase64), BASE64)
	//随机
	SetCall(Join(COLON, ModuleBuiltin, FuncNameRandomString), RandomString)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameRandomDigit), RandomDigit)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameRandom), Random)
	SetCall(Join(COLON, ModuleBuiltin, FuncNameUUID), UUID)
	//其他
	SetCall(Join(COLON, ModuleBuiltin, FuncNameLog), LOG)
}
