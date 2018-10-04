package builder

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	yaml "gopkg.in/yaml.v2"
)

// These types are defined in code.cloudfoundry.org/cli/util/manifest but importing it
// means importing 8 packages so redefined for simplicity instead.

type Manifest struct {
	Applications         []Application     `yaml:"applications"`
	Command              string            `yaml:"command,omitempty"`
	EnvironmentVariables map[string]string `yaml:"env,omitempty"`
}

type Application struct {
	Name                 string            `yaml:"name,omitempty"`
	Buildpack            string            `yaml:"buildpack,omitempty"`
	Buildpacks           []string          `yaml:"buildpacks,omitempty"`
	Command              string            `yaml:"command,omitempty"`
	EnvironmentVariables map[string]string `yaml:"env,omitempty"`
}

func readManifest(ctx context.Context, c client.Client) (*Manifest, error) {
	st := llb.Local(LocalNameContext,
		llb.SessionID(c.BuildOpts().SessionID),
		llb.IncludePatterns([]string{"manifest.yml"}),
		llb.SharedKeyHint("manifest.yml"),
		llb.WithCustomName("load manifest.yml"),
	)

	def, err := st.Marshal(llb.WithCaps(c.BuildOpts().LLBCaps))
	if err != nil {
		return nil, err
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, err
	}
	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}
	dt, err := ref.ReadFile(ctx, client.ReadRequest{
		Filename: "manifest.yml",
	})
	if err != nil {
		return nil, nil
	}

	var manifest Manifest
	if err := yaml.Unmarshal(dt, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
