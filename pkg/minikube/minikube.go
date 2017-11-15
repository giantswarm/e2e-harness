package minikube

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/giantswarm/e2e-harness/pkg/builder"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/micrologger"
)

type Minikube struct {
	logger   micrologger.Logger
	builder  builder.Builder
	imageTag string
}

func New(logger micrologger.Logger, builder builder.Builder, tag string) *Minikube {
	return &Minikube{
		logger:   logger,
		builder:  builder,
		imageTag: tag,
	}
}

// BuildImages is a Task that build the required images for both the main
// project and the e2e containers using the minikube docker environment.
func (m *Minikube) BuildImages() error {
	m.logger.Log("info", "Getting minikube docker environment")
	env, err := m.getDockerEnv()
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	name := harness.GetProjectName()
	image := fmt.Sprintf("quay.io/giantswarm/%s", name)
	m.logger.Log("info", "Building image "+image)
	if err := m.buildImage(name, dir, image, env); err != nil {
		return err
	}

	e2eBinary := name + "-e2e"
	e2eDir := filepath.Join(dir, "e2e")
	e2eImage := fmt.Sprintf("quay.io/giantswarm/%s-e2e", name)
	m.logger.Log("info", "Building image "+e2eImage)
	if err := m.buildImage(e2eBinary, e2eDir, e2eImage, env); err != nil {
		return err
	}

	return nil
}

func (m *Minikube) getDockerEnv() ([]string, error) {
	cmd := exec.Command("minikube", "docker-env")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []string{}, err
	}
	if err := cmd.Start(); err != nil {
		return []string{}, err
	}

	var env []string

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "export") {
			parts := strings.Fields(scanner.Text())
			entry := strings.Replace(parts[1], `"`, "", -1)
			env = append(env, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return env, err
	}

	if err := cmd.Wait(); err != nil {
		return []string{}, err
	}

	return env, nil
}

func compile(binaryName, path string) error {
	cmd := exec.Command("go", "build", "-o", binaryName, ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	cmd.Dir = path

	return cmd.Run()
}

func (m *Minikube) buildImage(binaryName, path, imageName string, env []string) error {
	if err := compile(binaryName, path); err != nil {
		fmt.Println("error compiling binary", binaryName)
		return err
	}

	if err := m.builder.Build(ioutil.Discard, imageName, path, m.imageTag, env); err != nil {
		fmt.Println("error building image", imageName)
		return err
	}
	return nil
}
