FROM golang:1.22-alpine3.20

ARG VERSION
ARG REVISION

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GO111MODULE=on

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
   sh scripts/build.sh

ENTRYPOINT ["/service/bin/main", "--config", "/etc/service/config.yaml"]
