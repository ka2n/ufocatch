#!/bin/sh

set -eu

BASE_DIR=$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)
cd $BASE_DIR

VERSION=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')
REPO="ufocatch"

rm -rf ./pkg
gox -os="linux darwin" -arch="386 amd64" -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

rm -rf ./dist
mkdir -p ${BASE_DIR}/dist/v${VERSION} 

for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    PLATFORM_NAME=$(basename ${PLATFORM})
    ARCHIVE_NAME=${REPO}_${VERSION}_${PLATFORM_NAME}
    pushd ./pkg/${PLATFORM_NAME}
    zip ${BASE_DIR}/dist/v${VERSION}/${ARCHIVE_NAME}.zip ./*
    popd
done
