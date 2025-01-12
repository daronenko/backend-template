.PHONY: docker/build docker/run docker/run/force docker/start docker/stop docker/down docker/clean docker/logs docker/logs/live docker/ps docker/ps/all docker/exec

DOCKER_COMPOSE_PATH ?= deploy/dev/docker-compose.yaml
DOCKER_COMPOSE 		?= docker compose -f ${DOCKER_COMPOSE_PATH} --profile ${PROFILE}

db: docker/build
docker/build:
	@VERSION=${VERSION} REVISION=${REVISION} ${DOCKER_COMPOSE} build

dr: docker/run
docker/run:
	@${DOCKER_COMPOSE} up -d

drf: docker/run/force
docker/run/force:
	@${DOCKER_COMPOSE} up -d --force-recreate

dst: docker/start
docker/start:
	@${DOCKER_COMPOSE} start ${c} ${compose}

dsp: docker/stop
docker/stop:
	@${DOCKER_COMPOSE} stop ${c} ${compose}

dd: docker/down
docker/down:
	@${DOCKER_COMPOSE} down

dc: docker/clean
docker/clean:
	@${DOCKER_COMPOSE} down -v --rmi all

dl: docker/logs
docker/logs:
	@${DOCKER_COMPOSE} logs ${c} ${compose}

dll: docker/logs/live
docker/logs/live:
	@${DOCKER_COMPOSE} logs -f ${c} ${compose}

dp: docker/ps
docker/ps:
	@${DOCKER_COMPOSE} ps

dpa: docker/ps/all
docker/ps/all:
	@${DOCKER_COMPOSE} ps -a

de: docker/exec
docker/exec:
	@${DOCKER_COMPOSE} exec -it ${c} ${compose} sh -c '(bash || ash || sh)'
