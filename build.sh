#!/usr/bin/env bash

build() {
  OS=${1}
  ARCH=${2}
  VERSION=${3}
  OUTPUT_DIR="intro-detection-info_${VERSION}_${OS}_${ARCH}"
  echo "Building ver. ${VERSION} for ${OS} (${ARCH}) to ${OUTPUT_DIR}..."
  mkdir ${OUTPUT_DIR}
  env GOOS=${OS} GOARCH=${ARCH} go build -ldflags="-w -s -X main.builtAt=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X main.version=${VERSION}" -o "${OUTPUT_DIR}/" github.com/jyggen/intro-detection-info
  cp CHANGELOG.md LICENSE README.md ${OUTPUT_DIR}/
  tar -czvf ${OUTPUT_DIR}.tar.gz ${OUTPUT_DIR}
  rm -rf ${OUTPUT_DIR}/
}

build darwin amd64 $1
build linux amd64 $1
build windows amd64 $1
