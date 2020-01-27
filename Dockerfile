# syntax=docker/dockerfile-upstream:experimental
FROM golang:1.12-alpine AS build
WORKDIR /go/src/github.com/tonistiigi/buildkit-pack
RUN apk add --no-cache file git
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
  GO111MODULE=on CGO_ENABLED=0 go build -o /out/pack ./cmd/pack && file /out/pack | grep "statically linked"
  
FROM scratch
COPY --from=build /out/pack /bin/pack
ENTRYPOINT ["/bin/pack"]
