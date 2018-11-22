package cmd

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/e2e-harness/cmd/internal"
	"github.com/giantswarm/e2e-harness/pkg/compiler"
	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/minikube"
	"github.com/giantswarm/e2e-harness/pkg/patterns"
	"github.com/giantswarm/e2e-harness/pkg/project"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/giantswarm/e2e-harness/pkg/wait"
)

var (
	TestCmd = &cobra.Command{
		Use:   "test",
		Short: "execute e2e tests",
		Run:   internal.NewRunFunc(runTest),
	}
)

var (
	testCmdTestDir string
)

func init() {
	RootCmd.AddCommand(TestCmd)

	TestCmd.Flags().StringVar(&testCmdTestDir, "test-dir", project.DefaultDirectory, "Name of the directory containing executable tests.")
}

func runTest(ctx context.Context, cmd *cobra.Command, args []string) error {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	projectTag := harness.GetProjectTag()
	projectName := harness.GetProjectName()

	fs := afero.NewOsFs()

	h := harness.New(logger, fs, harness.Config{})
	cfg, err := h.ReadConfig()
	if err != nil {
		return microerror.Mask(err)
	}

	// use latest tag for consumer projects (not dog-fooding e2e-harness)
	e2eHarnessTag := projectTag
	if projectName != "e2e-harness" {
		e2eHarnessTag = "latest"
	}

	var d *docker.Docker
	{
		c := docker.Config{
			Logger: logger,

			Dir:             testCmdTestDir,
			ExistingCluster: cfg.ExistingCluster,
			ImageTag:        e2eHarnessTag,
			RemoteCluster:   cfg.RemoteCluster,
		}

		d = docker.New(c)
	}

	pa := patterns.New(logger)
	w := wait.New(logger, d, pa)

	var p *project.Project
	{
		pCfg := &project.Config{
			Dir:  testCmdTestDir,
			Name: projectName,
			Tag:  projectTag,
		}
		pDeps := &project.Dependencies{
			Logger: logger,
			Runner: d,
			Wait:   w,
			Fs:     fs,
		}

		p = project.New(pDeps, pCfg)
	}

	var comp *compiler.Compiler
	{
		c := compiler.Config{
			Logger: logger,

			TestDir: testCmdTestDir,
		}

		comp = compiler.New(c)
	}

	// tasks to run
	bundle := []tasks.Task{
		comp.CompileTests,
	}

	if !cfg.RemoteCluster {
		// build images for minikube
		m := minikube.New(logger, d, projectTag)

		bundle = append(bundle, comp.CompileMain, m.BuildImages)
	}

	bundle = append(bundle, p.Test)

	return tasks.Run(ctx, bundle)
}
