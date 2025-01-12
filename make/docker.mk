.PHONY: docker/build docker/build/hot docker/run docker/down docker/clean docker/logs docker/logs/live

DOCKER_COMPOSE_PATH ?= deploy/dev/docker-compose.yaml
DOCKER_COMPOSE 		?= docker compose -f ${DOCKER_COMPOSE_PATH}

# compose (default: all)
c ?=

db: docker/build
docker/build:
	VERSION=${VERSION} REVISION=${REVISION} ${DOCKER_COMPOSE} --profile ${PROFILE} build

dr: docker/run
docker/run:
	@${DOCKER_COMPOSE} --profile ${PROFILE} up -d

dd: docker/down
docker/down:
	@${DOCKER_COMPOSE} --profile ${PROFILE} down

dc: docker/clean
docker/clean:
	@${DOCKER_COMPOSE} --profile ${PROFILE} down -v --rmi all

dl: docker/logs
docker/logs:
	@${DOCKER_COMPOSE} --profile ${PROFILE} logs ${c}

dll: docker/logs/live
docker/logs/live:
	@${DOCKER_COMPOSE} --profile ${PROFILE} logs -f ${c}
