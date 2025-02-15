include deploy/dev/.env

export VERSION := v0.0.1
export GIT_REVISION := ${shell git rev-parse HEAD}

include make/go.mk
include make/docker.mk
include make/swagger.mk
include make/migration.mk
