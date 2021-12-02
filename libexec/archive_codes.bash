#!/bin/bash

PROJECT=$1
shift
VERSION=$1
shift
TARGETDIR=$1
shift
ARCHIVEDIR=$1

SRTPATH=$(cd "$(dirname "$0")"; pwd)

if [ -z "$PROJECT" ]; then
    PROJECT=dragon
fi

if [ -z "$VERSION" ]; then
    VERSION=1.1.0
fi

if [ -z "$TARGETDIR" ]; then
    TARGETDIR=/vobs/cache/tmp
fi

if [ -z "$ARCHIVEDIR" ]; then
    ARCHIVEDIR=/vobs/tmp/
fi

cd $TARGETDIR

if [ -d $TARGETDIR/$PROJECT ]; then
    rm -rf $TARGETDIR/$PROJECT
fi

rsync -a $SRTPATH/../../dragon .

cd $TARGETDIR && tar czf ${PROJECT}_src_${VERSION}.tar.gz ${PROJECT}

cp ${PROJECT}_src_${VERSION}.tar.gz $ARCHIVEDIR

echo "${PROJECT} codes been archived in $ARCHIVEDIR${PROJECT}_src_${VERSION}.tar.gz"
