#!/bin/bash

set -ex

PROJECT_PATH=$(cd "$(dirname "${0}")"; pwd)
cd "$PROJECT_PATH"
mkdir -p bin

DRIVER_VERSION="1.0.0"
GIT_COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build -ldflags "-X github.com/opencurve/curvefs-csi/pkg/util.driverVersion=${DRIVER_VERSION} \
	-X github.com/opencurve/curvefs-csi/pkg/util.gitCommit=${GIT_COMMIT} \
	-X github.com/opencurve/curvefs-csi/pkg/util.buildDate=${BUILD_DATE}" \
	-o "${PROJECT_PATH}/bin/curvefs-csi-driver" "${PROJECT_PATH}/cmd/main.go"

echo "build successfully"
