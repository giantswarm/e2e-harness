package compiler

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/e2e-harness/pkg/harness"
)

type Config struct {
	Logger micrologger.Logger

	RemoteCluster bool
	TestDir       string
}

type Compiler struct {
	logger micrologger.Logger

	remoteCluster bool
	testDir       string
}

func New(config Config) *Compiler {
	c := &Compiler{
		logger: config.Logger,

		remoteCluster: config.RemoteCluster,
		testDir:       config.TestDir,
	}

	return c
}

// CompileMain is a Task that builds the main binary.
func (c *Compiler) CompileMain() error {
	dir, err := os.Getwd()
	if err != nil {
		return microerror.Mask(err)
	}

	mainPath := filepath.Join(dir, "main.go")
	_, err = os.Stat(mainPath)
	if os.IsNotExist(err) {
		c.logger.Log("function", "CompileMain", "level", "info", "message", "no main.go, skipping binary build")
		return nil
	}

	name := harness.GetProjectName()

	c.logger.Log("info", "Compiling binary "+name)
	if err := c.compileMain(name, dir); err != nil {
		c.logger.Log("info", "error compiling binary "+name)
		return microerror.Mask(err)
	}

	return nil
}

// CompileTests is a Task that builds the tests binary.
func (c *Compiler) CompileTests() error {
	dir, err := os.Getwd()
	if err != nil {
		return microerror.Mask(err)
	}

	e2eBinary := harness.GetProjectName() + "-e2e"
	e2eDir := filepath.Join(dir, c.testDir)

	c.logger.Log("info", "Compiling binary "+e2eBinary)
	err = c.compileTests(e2eBinary, e2eDir)
	if err != nil {
		c.logger.Log("info", "error compiling binary "+e2eBinary)
		return microerror.Mask(err)
	}

	return nil
}

// compileMain compiles a go binary in the given path giving it the provided
// name. If the binary already exists and is executable the build is skipped
func (c *Compiler) compileMain(binaryName, path string) error {
	// do not build if binary is already there
	binPath := filepath.Join(path, binaryName)
	if executebleExists(binPath) && c.remoteCluster {
		c.logger.Log("function", "compileMain", "level", "info", "message", "main binary exists, not building")
		return nil
	}

	cmd := exec.Command("go", "build", "-o", binaryName, ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux")
	cmd.Dir = path

	return cmd.Run()
}

// compileTests compiles a go test binary in the given path giving it the
// provided name. If the binary already exists and is executable the build
// is skipped
func (c *Compiler) compileTests(binaryName, path string) error {
	// do not build if binary is already there
	binPath := filepath.Join(path, binaryName)
	if executebleExists(binPath) && c.remoteCluster {
		c.logger.Log("function", "compileTests", "level", "info", "message", "test binary exists, not building")
		return nil
	}

	cmd := exec.Command("go", "test", "-c", "-o", binaryName, "-tags", "k8srequired", ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux")
	cmd.Dir = path
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func executebleExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	// 0111 octal represents a mode with the executable bit set
	// for user, group and others. Performing a bitwise and with
	// the mode of the file would only result in 0 if all these
	// bit flags are 0, so if the result is different from 0 we
	// can assume that the file is executable.
	isExecutable := fi.Mode()&os.FileMode(0111) != 0

	return isExecutable
}
