FROM golang:1.23-alpine3.20 AS builder

ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on

WORKDIR /app

ARG VERSION
ARG GIT_REVISION

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build \
    --ldflags "-X github.com/daronenko/backend-template/internal/pkg/bininfo.Version=${VERSION} -X github.com/daronenko/backend-template/internal/pkg/bininfo.GitRevision=${GIT_REVISION}" \
    -v -o ./bin/main ./cmd/app/main.go

FROM alpine:3.20

RUN apk add --no-cache curl

COPY --from=builder /app/bin/main /app

HEALTHCHECK --interval=10s --timeout=3s --retries=3 --start-interval=3s \
  CMD curl --fail http://localhost:8080/api/_/health || exit 1

ENTRYPOINT ["/app", "--config", "/etc/app/config.yaml"]
