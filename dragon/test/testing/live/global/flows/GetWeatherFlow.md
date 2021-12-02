************************************
# 流程
************************************
## 多次获取天气服务流程
> GetWeatherFlow
* m -> GetWeatherCtx.json k:GetWeatherCtx.json A:GetWeather.json B:GetWeather.json F:GetWeatherFlow2
// 无名或与x同名的上下文模板，使用系统内置的变量x，可以由外部覆盖取值
* A.a.nonce = x.nonce
* A.city = x.city
* F.nonce = 40000000 
// k所代表的上下文不会被覆盖，一般用于固定配置
* A.interval = k.interval
* B.a.nonce = 22222222
> GetWeatherFlow2
* m -> GetWeatherCtx.json A:GetWeather.json F:GetWeatherFlow3 B:GetWeather.json
* A.a.nonce = x.nonce
* F.nonce = 50000000
* B.a.nonce = 44444444
> GetWeatherFlow3
* m -> GetWeatherCtx.json A:GetWeather.json B:GetWeather.json
* A.a.nonce = x.nonce
* B.a.nonce = 55555555