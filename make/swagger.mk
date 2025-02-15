.PHONY: api

sa: api
api:
	@swag init -q -g cmd/app/main.go
