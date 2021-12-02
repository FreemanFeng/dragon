************************************
# 控制流
************************************
## 直播控制流
    插件 = P4
> C1 <- m
* Request m
* NowAddSeconds 3
* RandomDigit 3
************************************
# 前置条件
************************************
## 获取天气预报前置条件
> TS1 | scope = G
* Log 这是前置条件
************************************
# 后置条件
************************************
## 获取天气预报后置条件
> TT1 | scope = G
* Log 这是后置条件
************************************
# 用例
************************************
## 获取天气预报
    插件 = P4
> S1 <- C1 | setup = TS1 | teardown = TT1
### 调用模拟服务获取广州、上海各区天气预报
>> tc000010 | C = 1 | R = 1 | proto = http
*   // m -> A:2*Get*.json B:GetWeather.json
*   m -> A:GetWeather.json B:GetWeather.json
*  // A.a.ts = NowAddSeconds SECS
*   A.a.nonce = RandomDigit 8
*   B.a.nonce = RandomDigit 8
*   // A.a.sum = Add 3 9
*   // A.city = CITIES
*   // 1: A.district = GZ_DISTRICTS
*   // 2: A.district = SH_DISTRICTS
*   // A.interval = INTERVAL
*   // A -- city on A.a.nonce > 0
*   A.. -- long on A.a.nonce > 0
*   A...spots ++ name=soso id=888
*   A ++ hello=world
*   A.a.sign = Sign A.a.ts A.a.nonce 1 SECRET A.file
*   B.a.sign = Sign B.a.ts B.a.nonce 1 SECRET B.file
*   // => A.data[0] with message
*   //   A.data[1] without haha
*   // => A...message = B...message
*   => A.data = B.data
*   // => A.data.total = 5
*   //   A.data._total = 5
*   // =>  A...message with 暴雨 on city = 广州 district = 白云
*   // =>  A.code = 400 on A.city = "" or A.interval = 分钟
*   // =>  A.code = 400 on A.city = ""
*   // =>  A.code = 400
## 获取天气预报2
    插件 = P4
> S2 <- C1 | setup = TS1 | teardown = TT1
### 调用模拟服务获取广州、上海各区天气预报
>> tc000020 | C = 1 | R = 1 | proto = http
*   g -> B:GetWeather.json C:GetWeather.json
*   m -> GetWeatherCtx.json A:GetWeather.json F:GetWeatherFlow  2*g
*   A.a.nonce = 10000001
*   B.a.nonce = 10000002
*   C.a.nonce = 10000003
*   x.city = 广州
*   F.nonce = 20000000
*   A.a.sign = Sign A.a.ts A.a.nonce 1 SECRET A.file
*   B.city = A.city
*   =>  A...message with 暴雨 on city = 广州 district = 白云