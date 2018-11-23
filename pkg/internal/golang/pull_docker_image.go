package golang

import (
	"context"

	"github.com/giantswarm/e2e-harness/pkg/internal/docker"
	"github.com/giantswarm/microerror"
)

// PullDockerImage implements tasks.Task func type. It is meant to be ran
// before any other function from this package to add retries and avoid obscure
// pull logs in other tasks.
func PullDockerImage(ctx context.Context) error {
	err := docker.Pull(ctx, dockerImage)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
