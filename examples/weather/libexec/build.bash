#!/bin/bash
SRTPATH=$(cd "$(dirname "$0")" || exit; pwd)
#cd "$SRTPATH"/../src || exit
#go mod vendor
#cd "$SRTPATH"/../vendor || exit
#ln -sf ../deps/calc
cd "$SRTPATH"/../src || exit
go build -mod=vendor -o weather
mv weather "$SRTPATH"/../bin/

#cd "$SRTPATH"/../../../vendor || exit
#ln -sf ../testing/demo/plugins/src/demo
#ln -sf ../examples/weather/testing/weather/plugins/src/weather
