.PHONY: go/build go/run go/audit go/benchmark go/coverage go/format go/lint go/test go/longtest go/tidy

gb: go/build
go/build:
	@scripts/build.sh

gr: go/run
go/run: bin/main
	@bin/main

ga: go/audit
go/audit: check/govulncheck
	@go mod verify
	@go vet ./...
	@govulncheck ./...

gbm: go/benchmark
go/benchmark:
	@go test ./... -benchmem -bench=. -run=^Benchmark_$

gc: go/coverage
go/coverage: check/gotestsum
	@gotestsum -f testname -- ./... -race -count=1 -coverprofile=/tmp/coverage.out -covermode=atomic
	@go tool cover -html=/tmp/coverage.out

gf: go/format
go/format: check/gofumpt
	@go run mvdan.cc/gofumpt@latest -w -l .

gl: go/lint
go/lint: check/golangci-lint check/staticcheck
	@golangci-lint run ./...

gt: go/test
go/test: check/gotestsum
	@gotestsum -f testname -- ./... -race -count=1 -shuffle=on

glt: go/longtest
go/longtest: check/gotestsum
	@gotestsum -f testname -- ./... -race -count=15 -shuffle=on

gt: go/tidy
go/tidy:
	@go mod tidy -v
