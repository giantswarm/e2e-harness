package docker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/giantswarm/e2e-harness/pkg/harness"
)

func Run(imageTag string, out io.Writer, entrypoint string, args ...string) error {
	args = append([]string{"quay.io/giantswarm/e2e-harness:" + imageTag}, args...)

	dir, err := harness.BaseDir()
	if err != nil {
		return err
	}

	baseArgs := []string{
		"run",
		"-v", fmt.Sprintf("%s:%s", filepath.Join(dir, "workdir"), "/workdir"),
		"-e", fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", os.Getenv("AWS_ACCESS_KEY_ID")),
		"-e", fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", os.Getenv("AWS_SECRET_ACCESS_KEY")),
		"--entrypoint", entrypoint,
	}
	baseArgs = append(baseArgs, args...)

	cmd := exec.Command("docker", baseArgs...)
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd.Run()
}
