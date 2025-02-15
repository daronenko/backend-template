.PHONY: build run run/force enable disable shutdown clean logs logs/live ps ps/all shell

DOCKER_COMPOSE_PATH ?= deploy/dev/docker-compose.yaml
DOCKER_COMPOSE 		?= docker compose -f ${DOCKER_COMPOSE_PATH} --profile ${PROFILE}

db: build
build:
	@VERSION=${VERSION} REVISION=${REVISION} ${DOCKER_COMPOSE} build

dr: run
run:
	@${DOCKER_COMPOSE} up -d

drf: run/force
run/force:
	@${DOCKER_COMPOSE} up -d --force-recreate

de: enable
enable:
	@${DOCKER_COMPOSE} start ${c} ${compose}

dd: disable
disable:
	@${DOCKER_COMPOSE} stop ${c} ${compose}

ds: shutdown
shutdown:
	@${DOCKER_COMPOSE} down

dc: clean
clean:
	@${DOCKER_COMPOSE} down -v --rmi all

dl: logs
logs:
	@${DOCKER_COMPOSE} logs ${c} ${compose}

dll: logs/live
logs/live:
	@${DOCKER_COMPOSE} logs -f ${c} ${compose}

dp: ps
ps:
	@${DOCKER_COMPOSE} ps

dpa: ps/all
ps/all:
	@${DOCKER_COMPOSE} ps -a

dsh: shell
shell:
	@${DOCKER_COMPOSE} exec -it ${c} ${compose} sh -c '(bash || sh)'
