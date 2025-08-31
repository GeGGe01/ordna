# syntax=docker/dockerfile:1.7-labs

FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o /out/ordna ./cmd/ordna

FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=build /out/ordna /usr/local/bin/ordna
USER app
ENTRYPOINT ["/usr/local/bin/ordna"]
