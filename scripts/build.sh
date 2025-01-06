#!/bin/bash

set -euo pipefail

MAIN_PACKAGE="./cmd/app"
CONFIG_PACKAGE="github.com/daronenko/backend-template/internal/config"

OUTPUT_DIR="bin"
OUTPUT_FILE="${OUTPUT_DIR}/main"

: "${VERSION:?VERSION is not set. Please export VERSION and try again.}"
: "${REVISION:?REVISION is not set. Please export REVISION and try again.}"

mkdir -p "${OUTPUT_DIR}"

if [[ "${DEBUG:-0}" -eq "1" ]]; then
  echo "Building in debug mode..."
  go build \
    --gcflags="all=-N -l" \
    --ldflags "-X ${CONFIG_PACKAGE}.Version=${VERSION} -X ${CONFIG_PACKAGE}.Revision=${REVISION}" \
    -v -o "${OUTPUT_FILE}" "${MAIN_PACKAGE}"
else
  echo "Building in release mode..."
  go build \
    --ldflags "-X ${CONFIG_PACKAGE}.Version=${VERSION} -X ${CONFIG_PACKAGE}.Revision=${REVISION}" \
    -v -o "${OUTPUT_FILE}" "${MAIN_PACKAGE}"
fi

echo "Build complete."
