.PHONY: install audit benchmark coverage format lint test longtest vendor debug

gi: install
install:
	@go install mvdan.cc/gofumpt@latest
	@go install gotest.tools/gotestsum@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

ga: audit
audit:
	@go mod verify
	@go vet ./...
	@govulncheck ./...

gb: benchmark
benchmark:
	@go test ./... -benchmem -bench=. -run=^Benchmark_$

gc: coverage
coverage:
	@gotestsum -f testname -- ./... -race -count=1 -coverprofile=/tmp/coverage.out -covermode=atomic
	@go tool cover -html=/tmp/coverage.out

gf: format
format:
	@go run mvdan.cc/gofumpt@latest -w -l .

gl: lint
lint:
	@golangci-lint run ./...

gt: test
test:
	@gotestsum -f testname -- ./... -race -count=1 -shuffle=on

gl: longtest
longtest:
	@gotestsum -f testname -- ./... -race -count=15 -shuffle=on

gv: vendor
vendor:
	@go mod tidy -v
	@go mod vendor

gd: debug
debug:
	@dlv connect localhost:40000 
