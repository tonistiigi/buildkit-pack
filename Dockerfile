# syntax=tonistiigi/dockerfile:runmount20181002

FROM golang:1.11-alpine AS build
WORKDIR /go/src/github.com/tonistiigi/buildkit-pack
RUN apk add --no-cache file
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
  CGO_ENABLED=0 go build -o /out/pack ./cmd/pack && file /out/pack | grep "statically linked"
  
FROM scratch
COPY --from=build /out/pack /bin/pack
ENTRYPOINT ["/bin/pack"]