************************************
# 节点
************************************
## 获取天气
    使用内置变量x代表模板变量，用来替换
    路径 = /get/weather
    方法 = POST
    协议 = http
    主机 = HOST_WEATHER
> GetWeather
* a.nonce = RandomDigit 8
* interval = INTERVAL
* city = CITIES  
* a.sign = Sign a.ts a.nonce 1 SECRET file