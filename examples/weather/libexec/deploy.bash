#!/bin/bash

PLATFORM=$1
shift
PROJECT=$1

if [ -z "${PLATFORM}" ]; then
    PLATFORM=darwin
fi

if [ -z "${PROJECT}" ]; then
    PROJECT=weather
fi

SRTPATH=$(cd "$(dirname "$0")" || exit; pwd)


ROOT=${SRTPATH}/../.build/${PROJECT}
rm -rf "${ROOT}"/../*
mkdir -p "${ROOT}"

cd "${SRTPATH}"/../src || exit
CGO_ENABLED=0 GOOS=${PLATFORM} go build -a -ldflags '-extldflags "-static"' -o ${PROJECT}

mv ${PROJECT} "${ROOT}"/