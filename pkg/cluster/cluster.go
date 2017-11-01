package cluster

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/runner"
	"github.com/giantswarm/micrologger"
)

type Cluster struct {
	logger        micrologger.Logger
	runner        runner.Runner
	remoteCluster bool
}

func New(logger micrologger.Logger, runner runner.Runner, remoteCluster bool) *Cluster {
	return &Cluster{
		logger:        logger,
		runner:        runner,
		remoteCluster: remoteCluster,
	}
}

// Create is a Task that creates a remote cluster or, if we
// are using a local one, puts in place the required files for
// later access to it
func (c *Cluster) Create() error {
	if c.remoteCluster {
		return c.clusterAction("shipyard -action=start")
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = c.copyMinikubeAssets(usr.HomeDir)
	if err != nil {
		return err
	}
	err = c.setupMinikubeConfig(usr.HomeDir)
	if err != nil {
		return err
	}
	return nil
}

// Delete is a Task that gets rid of a remote cluster.
func (c *Cluster) Delete() error {
	return c.clusterAction("shipyard -action=stop")
}

func (c *Cluster) clusterAction(command string) error {
	if !c.remoteCluster {
		return nil
	}
	err := c.runner.Run(os.Stdout, command)

	return err
}

// copyMinikubeAssets copies all the files found in $HOME/.minikube to
// the e2e-harness workdir (so that they will be accessible from the test
// container)
func (c *Cluster) copyMinikubeAssets(homeDir string) error {
	c.logger.Log("info", "Making minikube assets accessible for the test container")

	originDir := filepath.Join(homeDir, ".minikube")
	baseDir, err := harness.BaseDir()
	if err != nil {
		return err
	}
	targetDir := filepath.Join(baseDir, "workdir", ".minikube")

	// copy minikube directory
	walkFn := func(path string, info os.FileInfo, err error) error {
		targetPath := strings.Replace(path, originDir, targetDir, 1)
		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}
		return copyFile(path, targetPath)
	}
	err = filepath.Walk(originDir, walkFn)
	if err != nil {
		return err
	}

	// copy kube config (assumes the current context is minukube)
	origKubeCfg := filepath.Join(homeDir, ".kube", "config")
	targetKubeCfg, err := getMinikubeConfigPath()
	if err != nil {
		return err
	}
	targetKubeCfgDir := filepath.Dir(targetKubeCfg)
	if err := os.MkdirAll(targetKubeCfgDir, os.ModePerm); err != nil {
		return err
	}
	if err := copyFile(origKubeCfg, targetKubeCfg); err != nil {
		return err
	}
	return nil
}

// setupMinikubeConfig replaces $HOME/.minukube in the k8s config
// file (as seen by the container where all the commands are going to
// be executed) by the path where the certificates can be found (again,
// from the container point of view).
func (c *Cluster) setupMinikubeConfig(homeDir string) error {
	c.logger.Log("info", "Setting up minikube config for the test container")

	// the default k8s config file references the required certificates
	// to access minikube using $HOME/.minikube, we store this in originDir
	originDir := filepath.Join(homeDir, ".minikube")

	// path is the actual location of the k8s config file that will be used from the
	// test container
	path, err := getMinikubeConfigPath()
	if err != nil {
		return err
	}
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// targetDir has the path where minikube certificates are stored as seen from
	// the test container
	targetDir := filepath.Join("/workdir", ".minikube")

	newContents := strings.Replace(string(read), originDir, targetDir, -1)

	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		return err
	}
	return nil
}

// getMinikubeConfigPath returns the actual path of the k8s config file that
// will be used by the test container (path from the point of view of the
// executing e2e-harness binary, not the test container).
func getMinikubeConfigPath() (string, error) {
	baseDir, err := harness.BaseDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(baseDir, "workdir", ".shipyard", "config")

	return path, nil
}

func copyFile(orig, dst string) error {
	in, err := os.Open(orig)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
