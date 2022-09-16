#!/usr/bin/env bash

build() {
  OS=${1}
  ARCH=${2}
  VERSION=${3}
  OUTPUT_DIR="intro-detection-info_${VERSION}_${OS}_${ARCH}"
  echo "Building ver. ${VERSION} for ${OS} (${ARCH}) to ${OUTPUT_DIR}..."
  mkdir ${OUTPUT_DIR}
  env GOOS=${OS} GOARCH=${ARCH} go build -trimpath -ldflags="-w -s -X main.builtAt=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X main.version=${VERSION}" -o "${OUTPUT_DIR}/" github.com/jyggen/intro-detection-info
  cp CHANGELOG.md LICENSE README.md ${OUTPUT_DIR}/

  if [ "${OS}" = "windows" ]; then
    zip -r ${OUTPUT_DIR}.zip ${OUTPUT_DIR}
  else
    tar -czvf ${OUTPUT_DIR}.tar.gz ${OUTPUT_DIR}
  fi
  rm -rf ${OUTPUT_DIR}/
}

build darwin amd64 $1
build darwin arm64 $1
build linux amd64 $1
build linux arm64 $1
build windows amd64 $1
