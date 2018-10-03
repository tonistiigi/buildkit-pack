package builder

import (
	"context"

	"github.com/moby/buildkit/frontend/gateway/client"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	return nil, errors.Errorf("not implemented")
}
