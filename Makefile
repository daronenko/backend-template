include deploy/dev/.env

VERSION ?= 0.0.1
export VERSION

GIT_COMMIT := $(shell git rev-parse HEAD)
export REVISION = $(GIT_COMMIT)

include make/go.mk
include make/goose.mk
include make/docker.mk
