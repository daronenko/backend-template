.PHONY: go/deps go/build go/run go/audit go/benchmark go/coverage go/format go/lint go/test go/longtest go/tidy

gd: go/deps
go/deps:
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install gotest.tools/gotestsum@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest

gb: go/build
go/build:
	@scripts/build.sh

gr: go/run
go/run: bin/main
	@bin/main

ga: go/audit
go/audit:
	@go mod verify
	@go vet ./...
	@govulncheck ./...

gbm: go/benchmark
go/benchmark:
	@go test ./... -benchmem -bench=. -run=^Benchmark_$

gc: go/coverage
go/coverage:
	@gotestsum -f testname -- ./... -race -count=1 -coverprofile=/tmp/coverage.out -covermode=atomic
	@go tool cover -html=/tmp/coverage.out

gf: go/format
go/format:
	@go run mvdan.cc/gofumpt@latest -w -l .

gl: go/lint
go/lint:
	@golangci-lint run ./...

gt: go/test
go/test:
	@gotestsum -f testname -- ./... -race -count=1 -shuffle=on

glt: go/longtest
go/longtest:
	@gotestsum -f testname -- ./... -race -count=15 -shuffle=on

gt: go/tidy
go/tidy:
	@go mod tidy -v

