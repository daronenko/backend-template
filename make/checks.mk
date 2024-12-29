.PHONY: check/staticcheck check/golangci-lint check/govulncheck check/gotestsum check/gofumpt check/goose

cs: check/staticcheck
check/staticcheck:
	@scripts/check_tool.sh staticcheck honnef.co/go/tools/cmd/staticcheck@latest

cl: check/golangci-lint
check/golangci-lint:
	@scripts/check_tool.sh golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint@latest

cv: check/govulncheck
check/govulncheck:
	@scripts/check_tool.sh govulncheck golang.org/x/vuln/cmd/govulncheck@latest

ct: check/gotestsum
check/gotestsum:
	@scripts/check_tool.sh gotestsum gotest.tools/gotestsum@latest

cf: check/gofumpt
check/gofumpt:
	@scripts/check_tool.sh gofumpt mvdan.cc/gofumpt@latest

cg: check/goose
check/goose:
	@scripts/check_tool.sh goose github.com/pressly/goose/v3/cmd/goose@latest
