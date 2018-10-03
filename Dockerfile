# syntax=tonistiigi/dockerfile:runmount20181002

FROM golang:1.11-alpine AS build
WORKDIR /go/src/github.com/tonistiigi/buildkit-pack
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
  go build -o /out/pack ./cmd/pack