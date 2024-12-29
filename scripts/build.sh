#!/bin/bash

set -euo pipefail

OUTPUT_DIR="bin"
MAIN_PACKAGE="./cmd"
OUTPUT_FILE="${OUTPUT_DIR}/main"

: "${VERSION:?VERSION is not set. Please export VERSION and try again.}"
: "${REVISION:?REVISION is not set. Please export REVISION and try again.}"

mkdir -p "${OUTPUT_DIR}"

if [[ "${DEBUG:-0}" -eq "1" ]]; then
  echo "Building in debug mode..."
  go build \
    --gcflags="all=-N -l" \
    --ldflags "-X main.Version=${VERSION} -X main.Revision=${REVISION}" \
    -v -o "${OUTPUT_FILE}" "${MAIN_PACKAGE}"
else
  echo "Building in release mode..."
  go build \
    --ldflags "-X main.Version=${VERSION} -X main.Revision=${REVISION}" \
    -v -o "${OUTPUT_FILE}" "${MAIN_PACKAGE}"
fi

echo "Build complete."
