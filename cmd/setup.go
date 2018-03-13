package cmd

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/e2e-harness/pkg/cluster"
	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/patterns"
	"github.com/giantswarm/e2e-harness/pkg/project"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/giantswarm/e2e-harness/pkg/wait"
)

var (
	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "setup e2e tests",
		RunE:  runSetup,
	}
)

var (
	setupCmdTestDir string
	name            string
	remoteCluster   bool
)

func init() {
	RootCmd.AddCommand(SetupCmd)

	SetupCmd.Flags().StringVar(&setupCmdTestDir, "test-dir", project.DefaultDirectory, "Name of the directory containing executable tests.")
	SetupCmd.Flags().StringVar(&name, "name", "e2e-harness", "CI execution identifier")
	SetupCmd.Flags().BoolVar(&remoteCluster, "remote", true, "use remote cluster")
}

func runSetup(cmd *cobra.Command, args []string) error {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	projectTag := harness.GetProjectTag()
	projectName := harness.GetProjectName()
	// use latest tag for consumer projects (not dog-fooding e2e-harness)
	e2eHarnessTag := projectTag
	if projectName != "e2e-harness" {
		e2eHarnessTag = "latest"
	}

	var d *docker.Docker
	{
		c := docker.Config{
			Logger: logger,

			Dir:           setupCmdTestDir,
			ImageTag:      e2eHarnessTag,
			RemoteCluster: remoteCluster,
		}

		d = docker.New(c)
	}

	pa := patterns.New(logger)
	w := wait.New(logger, d, pa)
	pCfg := &project.Config{
		Name: projectName,
		Tag:  projectTag,
	}
	fs := afero.NewOsFs()
	pDeps := &project.Dependencies{
		Logger: logger,
		Runner: d,
		Wait:   w,
		Fs:     fs,
	}
	p := project.New(pDeps, pCfg)
	hCfg := harness.Config{
		RemoteCluster: remoteCluster,
	}
	h := harness.New(logger, fs, hCfg)
	c := cluster.New(logger, fs, d, remoteCluster)

	// tasks to run
	bundle := []tasks.Task{
		h.Init,
		h.WriteConfig,
		c.Create,
		p.CommonSetupSteps,
	}

	return tasks.Run(bundle)
}
