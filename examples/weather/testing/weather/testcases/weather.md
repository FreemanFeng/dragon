************************************
# 控制流
************************************
## 获取天气预报控制流
> C1 <- m | plugins = P1
* Request m
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
> S1 <- C1 | plugins = P1 | setup = TS1 | teardown = TT1
### 调用模拟服务获取广州各区天气预报
* tc000010 | R = 5 | fill = random | proto = http
    m -> ClearNonce.c.json A:GetWeather.json B:GetWeather.json
    k = 0 1 6..8
    A.a.ts = NowAddSeconds k
    A.a.nonce = RandomDigit 8
    B.a.nonce = A.a.nonce B.a.ts = NowAddSeconds 0
    A.city = 广州 北京 深圳 上海 ""
    A.district = 天河 越秀 海珠 白云 on city = 广州
    A.district = 天河 on city = 北京
    A.interval = 分钟 on city = 深圳
    A -- locations city district on k = 1
    A.a.sign = Sign A.a.ts A.a.nonce 1 SECRET A.file
    
    => A.status_code = 400 on city = 北京 深圳 "" or k != 0
    A.status_code = 200 on city = 广州 上海
    A...message with 暴雨 on district = 白云
    B.status_code = 400