include .env
DB ?= postgres
DSN ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MIGRATIONS_PATH ?= db/migrations
SEED_PATH ?= db/seed

.PHONY: migration seed migrate/up migrate/down migrate/fresh migrate/fresh/seed

m: migration
migration: check/goose
	@goose -dir=${MIGRATIONS_PATH} create ${name} sql

s: seed
seed: check/goose
	@goose -dir=${SEED_PATH} create ${name} sql

mu: migrate/up
migrate/up: check/goose
	@goose -dir=${MIGRATIONS_PATH} ${DB} ${DSN} up

md: migrate/down
migrate/down: check/goose
	@goose -dir=${MIGRATIONS_PATH} ${DB} ${DSN} down

mf: migrate/fresh
migrate/fresh: check/goose
	@goose -dir=${MIGRATIONS_PATH} ${DB} ${DSN} reset
	@make migrate/up

msf: migrate/fresh/seed
migrate/fresh/seed: check/goose
	@make migrate/fresh
	@goose -dir=${MIGRATIONS_PATH} -no-versioning ${DB} ${DSN} up
