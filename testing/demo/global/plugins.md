************************************
# 插件
************************************
## 获取天气预报插件
    名称 = demo
> P1 | mode = rpc | port = 8088
## 演示调起bash脚本
    关闭 = 1
    名称 = add.bash
> P2 | mode = bin | calls = Add
## 获取天气预报插件
    关闭 = 1
    名称 = demo.py
    启动 = python3 demo.py 8086
> P3 | mode = src | port = 8086
## 获取天气预报插件
    关闭 = 1
> P4 | mode = rpc | port = 8081