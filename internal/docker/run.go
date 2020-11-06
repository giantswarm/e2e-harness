package docker

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/e2e-harness/v3/internal/exec"
)

type RunConfig struct {
	Volumes          []string
	Env              []string
	WorkingDirectory string
	Image            string
	Args             []string
}

func Run(ctx context.Context, config RunConfig) error {
	args := []string{
		"run",
		"--rm",
	}

	for _, volume := range config.Volumes {
		args = append(args, "-v", volume)
	}

	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	args = append(args, "-w", config.WorkingDirectory)
	args = append(args, config.Image)

	args = append(args, config.Args...)

	err := exec.Exec(ctx, "docker", args...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
