package cmd

import (
	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/patterns"
	"github.com/giantswarm/e2e-harness/pkg/project"
	"github.com/giantswarm/e2e-harness/pkg/results"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/giantswarm/e2e-harness/pkg/wait"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	TestCmd = &cobra.Command{
		Use:   "test",
		Short: "execute e2e tests",
		RunE:  runTest,
	}
)

func init() {
	RootCmd.AddCommand(TestCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	logger, err := micrologger.New(micrologger.DefaultConfig())
	if err != nil {
		return err
	}

	projectTag := harness.GetProjectTag()
	projectName := harness.GetProjectName()
	h := harness.New(logger, harness.Config{})
	cfg, err := h.ReadConfig()
	if err != nil {
		return err
	}

	// use latest tag for consumer projects (not dog-fooding e2e-harness)
	e2eHarnessTag := projectTag
	if projectName != "e2e-harness" {
		e2eHarnessTag = "latest"
	}

	d := docker.New(logger, e2eHarnessTag, cfg.RemoteCluster)
	pa := patterns.New(logger)
	w := wait.New(logger, d, pa)
	pCfg := &project.Config{
		Name: projectName,
		Tag:  projectTag,
	}
	fs := afero.NewOsFs()
	r := results.New(logger, fs, d)
	pDeps := &project.Dependencies{
		Logger:  logger,
		Runner:  d,
		Wait:    w,
		Results: r,
	}
	p := project.New(pDeps, pCfg)

	// tasks to run
	bundle := []tasks.Task{
		p.OutOfClusterTest,
		p.InClusterTest,
	}

	return tasks.Run(bundle)
}
