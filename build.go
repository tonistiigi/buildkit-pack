package builder

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	keyStack = "stack"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	opts := c.BuildOpts().Opts

	stack := "cflinuxfs2"
	if v, ok := opts[keyStack]; ok {
		stack = v
	}

	buildName, runName, err := builderImageName(stack)
	if err != nil {
		return nil, err
	}

	// TODO: read buildpacks and download directly

	// TODO: git/http sources
	src := llb.Local("context", llb.SessionID(c.BuildOpts().SessionID), llb.SharedKeyHint("pack-src"))

	builderImage := llb.Image(buildName, llb.WithMetaResolver(c))

	build := runBuilder(c, builderImage, `/packs/builder -buildpacksDir /var/lib/buildpacks  -outputDroplet /out/droplet.tgz -outputMetadata /out/result.json`, llb.Dir("/workspace"))
	build.AddMount("/workspace", src, llb.Readonly)
	build.AddMount("/tmp/cache", llb.Scratch(), llb.AsPersistentCacheDir("buildpack-build-cache", llb.CacheMountShared))

	extract := llb.Image("alpine").Run(llb.Shlex(`sh -c "mkdir -p /out/home/vcap && tar -C /out/home/vcap -xzf /in/droplet.tgz && chown -R 2000:2000 /out/home/vcap"`), llb.WithCustomName("copy droplet to stack"), llb.Dir("/in"))

	extract.AddMount("/in", build.Root(), llb.SourcePath("out"), llb.Readonly)
	st := extract.AddMount("/out", llb.Image(runName))

	def, err := st.Marshal()
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(ctx)

	var res *client.Result
	eg.Go(func() error {
		r, err := c.Solve(ctx, client.SolveRequest{
			Definition: def.ToPB(),
		})
		if err != nil {
			return err
		}
		res = r
		return nil
	})

	var config []byte
	eg.Go(func() error {
		_, c, err := c.ResolveImageConfig(ctx, runName, client.ResolveImageConfigOpt{})
		if err != nil {
			return err
		}
		config = c
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	// TODO: is the build label needed?
	res.AddMeta(exptypes.ExporterImageConfigKey, config)

	return res, nil
}

func runBuilder(c client.Client, img llb.State, cmd string, opts ...llb.RunOption) llb.ExecState {
	// work around docker 18.06 executor with no cgroups mounted because build has
	// a hard requirement on the file

	caps := c.BuildOpts().LLBCaps

	mountCgroups := (&caps).Supports(pb.CapExecCgroupsMounted) != nil

	opts = append(opts, llb.WithCustomName(cmd))

	if mountCgroups {
		cmd = `sh -c "mkdir -p /sys/fs/cgroup/memory && echo 9223372036854771712 > /sys/fs/cgroup/memory/memory.limit_in_bytes && ` + cmd + `"`
	}

	es := img.Run(append(opts, llb.Shlex(cmd))...)

	if mountCgroups {
		es.AddMount("/sys/fs/cgroup", llb.Scratch())
		hosts := llb.Image("alpine").Run(llb.Shlex(`sh -c 'echo "127.0.0.1 $(hostname)" > /out/hosts'`), llb.WithCustomName("[internal] make hostname resolvable"))
		es.AddMount("/etc/hosts", hosts.Root(), llb.SourcePath("hosts"), llb.Readonly)
	}

	return es
}

func builderImageName(stack string) (string, string, error) {
	switch stack {
	case "cflinuxfs2":
		return "docker.io/packs/cflinuxfs2:build", "docker.io/packs/cflinuxfs2:run", nil
	default:
		return "", "", errors.Errorf("unsupported stack %s", stack)
	}
}
