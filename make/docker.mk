.PHONY: docker/build docker/build/hot docker/run docker/down docker/clean docker/logs docker/logs/live

DOCKER_COMPOSE_PATH ?= deploy/dev/docker-compose.yaml
DOCKER_COMPOSE 		?= docker compose -f ${DOCKER_COMPOSE_PATH}

db: docker/build
docker/build:
	@VERSION=${VERSION} REVISION=${REVISION} ${DOCKER_COMPOSE} build

dr: docker/run
docker/run:
	@${DOCKER_COMPOSE} up -d

dd: docker/down
docker/down:
	@${DOCKER_COMPOSE} down

dc: docker/clean
docker/clean:
	@${DOCKER_COMPOSE} down -v --rmi all

c       ?=
compose ?=

dl: docker/logs
docker/logs:
	@${DOCKER_COMPOSE} logs ${compose} ${c}

c       ?=
compose ?=

dll: docker/logs/live
docker/logs/live:
	@${DOCKER_COMPOSE} logs -f ${compose} ${c}
