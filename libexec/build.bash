#!/bin/bash
SRTPATH=$(cd "$(dirname "$0")" || exit; pwd)
#cd "$SRTPATH"/.. || exit
#go mod init
cd "$SRTPATH"/../dragon || exit
go mod vendor
cd "$SRTPATH"/../vendor || exit
ln -sf ../deps/calc
cd "$SRTPATH"/../dragon || exit
go build -mod=vendor

cd "$SRTPATH"/../vendor || exit
ln -sf ../testing/demo/plugins/src/demo
ln -sf ../examples
ln -sf ../examples/weather/testing/weather/plugins/src/weather
ln -sf ../goplugin
