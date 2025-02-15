.PHONY: migration migrate/up migrate/down migrate/fresh seed seed/fresh

DB 	?= postgres
DSN ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSLMODE}

MIGRATIONS_PATH ?= db/migrations
SEEDS_PATH      ?= db/seed

MIGRATIONS ?= goose -dir=${MIGRATIONS_PATH}
SEEDS      ?= goose -dir=${SEEDS_PATH}

name ?= unknown

m: migration
migration:
	@${MIGRATIONS} create ${name} sql

mu: migrate/up
migrate/up:
	@${MIGRATIONS} ${DB} ${DSN} up

md: migrate/down
migrate/down:
	@${MIGRATIONS} ${DB} ${DSN} down

mf: migrate/fresh
migrate/fresh:
	@${MIGRATIONS} ${DB} ${DSN} reset
	@make migrate/up

s: seed
seed:
	@${SEEDS} create ${name} sql

sf: seed/fresh
seed/fresh:
	@make migrate/fresh
	@${MIGRATIONS} -no-versioning ${DB} ${DSN} up
