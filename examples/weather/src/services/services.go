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

package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/nuxim/dragon/examples/weather/src/common"
)

type GlobalConfig struct {
	Nonce sync.Map // 随机值
}

var mG *GlobalConfig
var oG sync.Once

func GetGlobalConfig() *GlobalConfig {
	oG.Do(func() {
		mG = &GlobalConfig{Nonce: sync.Map{}}
	})
	return mG
}

func ExistsNonce(key int64) bool {
	k := GetGlobalConfig()
	_, ok := k.Nonce.Load(key)
	return ok
}

func SetNonce(key int64) {
	k := GetGlobalConfig()
	k.Nonce.Store(key, 1)
}

func GetWeather(ts, nonce int64, sign string, msg *WeatherRequest) WeatherResponse {
	rand.Seed(time.Now().Unix())
	fmt.Printf("ts:%d nonce:%d sign:%s msg:%s", ts, nonce, sign, msg)
	rsp := WeatherResponse{}
	rsp.Echo = nonce
	data := WeatherData{Time: "2021.05.09", City: msg.City, District: msg.District, Message: "天晴"}
	rsp.Data = append(rsp.Data, data, data, data, data, data)
	rsp.Data[1].Time = "2021.05.10"
	rsp.Data[1].Message = "暴雨"
	rsp.Data[1].City = "上海"
	rsp.Data[1].District = "长宁"
	rsp.Data[2].Time = "2021.05.11"
	rsp.Data[2].Message = "阴天"
	rsp.Data[2].City = "深圳"
	rsp.Data[2].District = "南山"
	rsp.Data[3].Time = "2021.05.12"
	rsp.Data[3].Message = "下雪"
	rsp.Data[3].City = "北京"
	rsp.Data[3].District = "朝阳"
	rsp.Data[4].Time = "2021.05.13"
	rsp.Data[4].Message = "闪电"
	rsp.Data[4].City = "上海"
	rsp.Data[4].District = "闵行"
	if msg.City == "广州" && msg.District == "白云" {
		for i := 0; i < 5; i++ {
			s := []string{"暴雨", strconv.Itoa(i)}
			rsp.Data[i].Message = strings.Join(s, "")
		}
	}
	fmt.Println("-->应答数据:", rsp)
	SetNonce(nonce)
	return rsp
}
