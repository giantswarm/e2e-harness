package docker

import (
	"context"
	"fmt"

	"github.com/giantswarm/e2e-harness/pkg/internal/exec"
	"github.com/giantswarm/microerror"
)

type RunConfig struct {
	Rm               bool
	Volumes          []string
	Env              []string
	WorkingDirectory string
	Image            string
	Args             []string
}

func Run(ctx context.Context, config RunConfig) error {
	prog := "docker"
	args := []string{
		"run",
		fmt.Sprintf("--rm=%t", config.Rm),
	}

	for _, volume := range config.Volumes {
		args = append(args, "-v", volume)
	}

	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	args = append(args, "-w", config.WorkingDirectory)
	args = append(args, config.Image)

	for _, arg := range config.Args {
		args = append(args, arg)
	}

	err := exec.Exec(ctx, prog, args...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
