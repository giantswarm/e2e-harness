package compiler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/internal/golang"
)

type Config struct {
	Logger micrologger.Logger

	TestDir string
}

type Compiler struct {
	logger micrologger.Logger

	testDir string
}

func New(config Config) *Compiler {
	c := &Compiler{
		logger: config.Logger,

		testDir: config.TestDir,
	}

	return c
}

// CompileMain is a Task that builds the main binary.
func (c *Compiler) CompileMain(ctx context.Context) error {
	binaryPath := harness.GetProjectName()

	c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("compiling binary %#q", binaryPath))

	dir, err := os.Getwd()
	if err != nil {
		return microerror.Mask(err)
	}

	mainPath := filepath.Join(dir, "main.go")
	_, err = os.Stat(mainPath)
	if os.IsNotExist(err) {
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("did not compile binary %#q", binaryPath))
		c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("file main.go was not found"))
		return nil
	}

	err = golang.Go(ctx, "build", "-o", binaryPath, ".")
	if err != nil {
		return microerror.Mask(err)
	}

	c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("compiled binary %#q", binaryPath))
	return nil
}

// CompileTests is a Task that builds the tests binary.
func (c *Compiler) CompileTests(ctx context.Context) error {
	binaryPath := filepath.Join(c.testDir, harness.GetProjectName()+"-e2e")

	c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("compiling binary %#q", binaryPath))

	err := golang.Go(ctx, "test", "-c", "-o", binaryPath, "-tags", "k8srequired", c.testDir)
	if err != nil {
		return microerror.Mask(err)
	}

	c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("compiled binary %#q", binaryPath))
	return nil
}
