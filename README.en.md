# dragon

#### Description
Rules Driven Simplified and Efficient Automation Test Framework - Dragon

#### Software Architecture


#### Installation
go version: 1.17+
> Linux configure GO proxy
* export GOPROXY=https://goproxy.io
* export GO111MODULE=on
> Windows configure GO proxy
* go env -w GOPROXY="https://goproxy.io"
* go env -w GO111MODULE="on"
> Installation Steps
1.  git clone git@gitee.com:freemanfeng/dragon.git
2.  cd dragon
3.  go mod vendor
4.  cd vendor
5.  ln -sf ../goplugin
6.  ln -sf ../testing/demo/plugins/src/demo

#### Instructions

1.  git clone git@gitee.com:freemanfeng/dragon.git
2.  cd dragon/dragon
3.  go build
4.  ./dragon

#### Demo
##### Open two Linux Terminals，buildup Demo and Weather services
> Demo service
* cd dragon/testing/demo/plugins/src/demo
* go build
* ./demo

> Weather service
* cd dragon/examples/weather/src
* go build
* mv src weather
* ./weather
##### Open two Windows cmd，buildup Demo and Weather services
> Demo service
* cd dragon/testing/demo/plugins/src/demo
* go build
* demo.exe

> Weather service
* cd dragon/examples/weather/src
* go build
* mv src.exe weather.exe
* weather.exe

at last, open browser and visit http://localhost:9899/test/demo?c=tc000020