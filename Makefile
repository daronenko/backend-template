include deploy/dev/.env

export VERSION := 0.0.1

GIT_COMMIT := $(shell git rev-parse HEAD)
export REVISION := $(GIT_COMMIT)

include make/go.mk
include make/goose.mk
include make/docker.mk
