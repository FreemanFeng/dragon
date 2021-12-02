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

type WeatherRequest struct {
	Interval string   // 间隔类型，小时/天
	City     string   // 城市
	District string   // 区
	Good     GoodType // 好去处
	Hello    string
}

type GoodType struct {
	Spots SpotsType // 景点
}

type SpotsType struct {
	Pos  []Location // 位置列表
	Name string
	Id   int
}

type Location struct {
	Long string // 经度
	Lat  string // 维度
}

type WeatherResponse struct {
	Echo int64         `json:"echo"` // 防重放
	Data []WeatherData `json:"data"` // 天气数据
}

type WeatherData struct {
	Time     string `json:"time"`     // 时间
	City     string `json:"city"`     // 城市
	District string `json:"district"` // 区
	Message  string `json:"message"`  // 天气提醒
}
