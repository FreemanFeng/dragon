#!/bin/bash
VER=$1
shift
PROJECT=$1

if [ -z "${VER}" ]; then
    VER=0.1.0
fi

if [ -z "${PROJECT}" ]; then
    PROJECT=dragon
fi


SRTPATH=$(cd "$(dirname "$0")" || exit; pwd)


ROOT=${SRTPATH}/../.build/${PROJECT}
rm -rf "${ROOT}"
mkdir -p "${ROOT}"
mkdir -p "${ROOT}"/testing

for platform in linux darwin windows; do
  cd "${SRTPATH}"/../dragon || exit
  rm -rf "${ROOT:?}"/${PROJECT}*
  out=${PROJECT}
  if [ "$platform" == "windows" ]; then
    out=${PROJECT}.exe
  fi
#  CGO_ENABLED=0 GOOS=${platform} GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o ${out}
  GOOS=${platform} GOARCH=amd64 go build -mod=vendor -o ${out}
  mv ${out} "${ROOT}"/
  rsync -a "${SRTPATH}"/../go.mod "${ROOT}"/

  cd "${SRTPATH}"/../examples/weather/src || exit
  rm -rf "${ROOT}"/weather*
  out=weather
  if [ "$platform" == "windows" ]; then
    out=weather.exe
  fi
  CGO_ENABLED=0 GOOS=${platform} GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o ${out}
  mv ${out} "${ROOT}"/

  cd "${SRTPATH}"/../testing || exit
  find . -name \*.so -exec rm -rf {} \;
  cd "${SRTPATH}"/../testing/demo/plugins/src || exit
  rsync -a "${SRTPATH}"/../goplugin demo/
  rsync -a "${SRTPATH}"/../pyplugin pydemo/
  cd "${SRTPATH}"/../testing || exit
  rsync -a demo "${ROOT}"/testing/
  find "${ROOT}"/testing/ -name \*.exe -exec rm -rf {} \;
  find "${ROOT}"/testing/ -name \*.so -exec rm -rf {} \;
  find "${ROOT}"/testing/ -name .DS_Store -exec rm -rf {} \;
  find "${ROOT}"/testing/ -name .\* -exec rm -rf {} \;

  cd "${ROOT}"/.. || exit

  if [ "$platform" == "windows" ]; then
    zip -r -q "${PROJECT}_v${VER}_${platform}.zip" ${PROJECT}
    echo "Archived bin to ${ROOT}/../${PROJECT}_v${VER}_${platform}.zip"
  else
    tar czf "${PROJECT}_v${VER}_${platform}.tar.gz" ${PROJECT}
    echo "Archived bin to ${ROOT}/../${PROJECT}_v${VER}_${platform}.tar.gz"
  fi

done
