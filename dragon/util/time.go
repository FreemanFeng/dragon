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

package util

import (
	"strconv"
	"strings"
	"time"

	. "github.com/FreemanFeng/dragon/dragon/common"
)

func ConvertTimeDigit(date time.Time, days int) string {
	ret := date.AddDate(0, 0, days)
	return ret.Format("20060102150405")
}

func GetTimeDigit(days int) string {
	now := time.Now().Local()
	return ConvertTimeDigit(now, days)
}

func Timestamp(days int) int64 {
	now := time.Now().Local()
	ts := now.AddDate(0, 0, days)
	return ts.Unix()
}

func Millisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Seconds() int64 {
	return time.Now().Unix()
}

func CalculateSecs(s string) int64 {
	n := len(s)
	if n == 0 {
		return 0
	}
	k, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		k, e = strconv.ParseInt(s[:n-1], 10, 64)
		if e != nil {
			return 0
		}
	}
	if StrMatchEnd(s, "m", "M") {
		return k * 60
	}
	if StrMatchEnd(s, "h", "H") {
		return k * 3600
	}
	if StrMatchEnd(s, "d", "D") {
		return k * 3600 * 24
	}
	return 0
}

func ConvertTime(date time.Time, days int, format string) string {
	ret := date.AddDate(0, 0, days)
	if format == EMPTY {
		return ret.Format("2006-01-02 15:04:05")
	}
	return ret.Format(format)
}

func ConvertDate(date time.Time, days int, format string) string {
	ret := date.AddDate(0, 0, days)
	if format == EMPTY {
		return ret.Format("2006-01-02")
	}
	return ret.Format(format)
}

func GetTime(days int) string {
	return GetTimeFormat(days, EMPTY)
}

func GetTimeFormat(days int, format string) string {
	now := time.Now().Local()
	return ConvertTime(now, days, format)
}

func GetDate(days int) string {
	return GetDateFormat(days, EMPTY)
}

func GetDateFormat(days int, format string) string {
	now := time.Now().Local()
	return ConvertDate(now, days, format)
}

func GetHour() int { //获取当前服务器的小时
	unixTimeStamp := time.Now().Unix()
	timeStamp := time.Unix(unixTimeStamp, 0)
	hr, _, _ := timeStamp.Clock()
	return hr
}

func ConvertMonth(date string) string {
	s := strings.Split(date, "-")
	m := map[string]string{"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr", "05": "May",
		"06": "Jun", "07": "Jul", "08": "Aug", "09": "Sep",
		"1": "Jan", "2": "Feb", "3": "Mar", "4": "Apr", "5": "May",
		"6": "Jun", "7": "Jul", "8": "Aug", "9": "Sep", "10": "Oct", "11": "Nov", "12": "Dec"}
	d := s[1]
	h := []string{s[0], m[d], s[2]}
	return strings.Join(h, "-")
}

func DateOffset(date string) int {
	s := ConvertMonth(date)
	loc, _ := time.LoadLocation("Local")
	d, _ := time.ParseInLocation("2006-Jan-02", s, loc)
	now := time.Now().Local().Unix()
	return int(d.Unix()-now) / DaySecs
}

func DateOffsetString(date string) string {
	offset := DateOffset(date)
	return strconv.Itoa(offset)
}

func ConvertSeconds(seconds int64, days int) string {
	return ConvertSecondsFormat(seconds, days, EMPTY)
}

func ConvertSecondsFormat(secs int64, days int, format string) string {
	t := time.Unix(Seconds(), 0)
	if secs != UNDEFINED {
		t = time.Unix(secs, 0)
	}
	return ConvertTime(t, days, format)
}

func GetYesterday(date string) string {
	return GetYesterdayFormat(date, EMPTY)
}

func GetYesterdayFormat(date string, format string) string {
	s := ConvertMonth(date)
	loc, _ := time.LoadLocation("Local")
	d, _ := time.ParseInLocation("2006-Jan-02", s, loc)
	return ConvertDate(d, -1, format)
}

func GetLastWeek(date string) string {
	return GetLastWeekFormat(date, EMPTY)
}

func GetLastWeekFormat(date string, format string) string {
	s := ConvertMonth(date)
	loc, _ := time.LoadLocation("Local")
	d, _ := time.ParseInLocation("2006-Jan-02", s, loc)
	return ConvertDate(d, -7, format)
}

func GetCustomDays(date string, cs int) string {
	return GetCustomDaysFormat(date, cs, EMPTY)
}

func GetCustomDaysFormat(date string, cs int, format string) string {
	s := ConvertMonth(date)
	loc, _ := time.LoadLocation("Local")
	d, _ := time.ParseInLocation("2006-Jan-02", s, loc)
	return ConvertDate(d, cs, format)
}

func GetLastMonth(date string) string {
	return GetLastMonthFormat(date, EMPTY)
}

func GetLastMonthFormat(date, format string) string {
	s := ConvertMonth(date)
	loc, _ := time.LoadLocation("Local")
	d, _ := time.ParseInLocation("2006-Jan-02", s, loc)
	return ConvertDate(d, -30, format)
}

func GetTimestamp(days int) string {
	ts := Timestamp(days)
	return strconv.FormatInt(ts, 10)
}

func GetMillisecond() string {
	ts := Millisecond()
	return strconv.FormatInt(ts, 10)
}

func GetSeconds() string {
	ts := Seconds()
	return strconv.FormatInt(ts, 10)
}

func AddSeconds(secs int64) string {
	ts := Seconds()
	return strconv.FormatInt(ts+secs, 10)
}
