package docker

import (
	"context"

	"github.com/giantswarm/e2e-harness/pkg/internal/exec"
	"github.com/giantswarm/microerror"
)

func Pull(ctx context.Context, image string) error {
	args := []string{
		"pull",
		image,
	}

	err := exec.Exec(ctx, "docker", args...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
