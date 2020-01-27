module github.com/tonistiigi/buildkit-pack

go 1.12

require (
	github.com/Microsoft/go-winio v0.4.14
	github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50
	github.com/containerd/containerd v1.4.0-0.20191014053712-acdcf13d5eaf
	github.com/containerd/continuity v0.0.0-20200107194136-26c1120b8d41
	github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/google/shlex v0.0.0-20150127133951-6f45313302b9
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/moby/buildkit v0.0.0-20181003224033-f07efb78e3f1
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opentracing/opentracing-go v0.0.0-20171003133519-1361b9cd60be
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	github.com/tonistiigi/fsutil v0.0.0-20191018213012-0f039a052ca1
	github.com/tonistiigi/units v0.0.0-20180711220420-6950e57a87ea
	github.com/urfave/cli v0.0.0-20171014202726-7bc6a0acffa5
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190812073006-9eafafc0a87e
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20161028155119-f51c12702a4d
	google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
	google.golang.org/grpc v1.23.0
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/moby/buildkit => github.com/hinshun/buildkit v0.0.0-20200124224350-99d18890d310

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
