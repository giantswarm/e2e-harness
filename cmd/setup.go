package cmd

import (
	"github.com/giantswarm/e2e-harness/pkg/cluster"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
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
	// tasks to run
	bundle := []tasks.Task{
		harness.Init,
		cluster.Create,
		harness.WriteStatus,
	}

	status := harness.Status{
		RemoteCluster: remoteCluster,
		GitCommit:     GetGitCommit(),
	}

	return tasks.Run(bundle, status)
}
