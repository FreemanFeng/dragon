# dragon

#### 介绍
简洁高效规则驱动自动化测试框架Dragon

#### 软件架构
软件架构说明


#### 安装教程
开发环境的go版本1.17+
> Linux配置GO代理
* export GOPROXY=https://goproxy.io
* export GO111MODULE=on
> Windows配置GO代理
* go env -w GOPROXY="https://goproxy.io"
* go env -w GO111MODULE="on"
> 安装步骤
1.  git clone git@gitee.com:freemanfeng/dragon.git
2.  cd dragon
3.  go mod vendor
4.  cd vendor
5.  ln -sf ../goplugin
6.  ln -sf ../testing/demo/plugins/src/demo

#### 使用说明

1.  git clone git@gitee.com:freemanfeng/dragon.git
2.  cd dragon/dragon
3.  go build
4.  ./dragon

#### 演示
##### Linux下打开两个Terminal，分别编译和启动Demo和Weather服务
> Demo服务
* cd dragon/testing/demo/plugins/src/demo
* go build
* ./demo

> Weather服务
* cd dragon/examples/weather/src
* go build
* mv src weather  
* ./weather
##### Windows下打开两个cmd，分别编译和启动Demo和Weather服务
> Demo服务
* cd dragon/testing/demo/plugins/src/demo
* go build
* demo.exe

> Weather服务
* cd dragon/examples/weather/src
* go build
* mv src.exe weather.exe
* weather.exe

最后，在浏览器上访问 http://localhost:9899/test/demo?c=tc000020