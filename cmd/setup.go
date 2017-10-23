package cmd

import (
	"github.com/giantswarm/e2e-harness/pkg/cluster"
	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/project"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
)

var (
	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "setup e2e tests",
		RunE:  runSetup,
	}
	remoteCluster bool
)

func init() {
	RootCmd.AddCommand(SetupCmd)

	SetupCmd.Flags().BoolVar(&remoteCluster, "remote-cluster", true, "use remote cluster")
}

func runSetup(cmd *cobra.Command, args []string) error {
	logger, err := micrologger.New(micrologger.DefaultConfig())
	if err != nil {
		return err
	}

	gitCommit := GetGitCommit()

	d := docker.New(logger, gitCommit)
	p := project.New(logger, d)
	hCfg := harness.Config{
		RemoteCluster: remoteCluster,
	}
	h := harness.New(logger, hCfg)
	c := cluster.New(logger, d, remoteCluster)

	// tasks to run
	bundle := []tasks.Task{
		h.Init,
		c.Create,
		p.CommonSetupSteps,
		p.SetupSteps,
		h.WriteConfig,
	}

	return tasks.Run(bundle)
}
