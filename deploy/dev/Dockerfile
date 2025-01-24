FROM golang:1.23-alpine3.20 AS builder

ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on

WORKDIR /service

ARG VERSION
ARG REVISION

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build \
    --ldflags "-X github.com/daronenko/backend-template/internal/config.Version=${VERSION} -X github.com/daronenko/backend-template/internal/config.Revision=${REVISION}" \
    -v -o ./bin/main ./cmd/app/main.go

FROM scratch

COPY --from=builder /service/bin/main /service

ENTRYPOINT ["/service", "--config", "/etc/service/config.yaml"]
