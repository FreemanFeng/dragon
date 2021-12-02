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
	"strconv"

	. "github.com/nuxim/dragon/dragon/common"
	"github.com/nuxim/dragon/dragon/util"
)

func getFormatParams(params ...string) (int, string, int64) {
	n := len(params)
	seconds := int64(UNDEFINED)
	days := 0
	format := EMPTY
	if n > 0 {
		days, _ = strconv.Atoi(params[0])
	}
	if n > 1 {
		format = params[1]
	}
	if n > 2 {
		seconds, _ = strconv.ParseInt(params[2], 10, 64)
	}
	return days, format, seconds
}

func FormatDate(params ...string) string {
	days, format, _ := getFormatParams(params...)
	return util.GetDateFormat(days, format)
}

func DATE(params ...string) string {
	return FormatDate(params...)
}

func FormatTime(params ...string) string {
	days, format, _ := getFormatParams(params...)
	return util.GetTimeFormat(days, format)
}

func TIME(params ...string) string {
	return FormatTime(params...)
}

func FormatSeconds(params ...string) string {
	days, format, secs := getFormatParams(params...)
	return util.ConvertSecondsFormat(secs, days, format)
}

func SEC(params ...string) string {
	return FormatSeconds(params...)
}

func FormatNowInSeconds(params ...string) string {
	return util.GetSeconds()
}

func NOW(params ...string) string {
	return util.GetSeconds()
}

func NowAddSeconds(params ...string) string {
	secs := int64(0)
	if len(params) > 0 {
		k, err := strconv.ParseInt(params[0], 10, 64)
		if err == nil {
			secs = k
		}
	}
	return util.AddSeconds(secs)
}

func NowAddDays(params ...string) string {
	days := 0
	if len(params) > 0 {
		k, err := strconv.Atoi(params[0])
		if err == nil {
			days = k
		}
	}
	return util.GetDate(days)
}

func NowAddHours(params ...string) string {
	hours := 0
	if len(params) > 0 {
		k, err := strconv.Atoi(params[0])
		if err == nil {
			hours = k
		}
	}
	t := util.GetHour() + hours
	return strconv.Itoa(t)
}

func NowMS(params ...string) string {
	return util.GetMillisecond()
}
