package cmd

import (
	"github.com/giantswarm/e2e-harness/pkg/cluster"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/spf13/cobra"
)

var (
	TeardownCmd = &cobra.Command{
		Use:   "teardown",
		Short: "teardown e2e tests",
		RunE:  runTeardown,
	}
)

func init() {
	RootCmd.AddCommand(TeardownCmd)
}

func runTeardown(cmd *cobra.Command, args []string) error {
	bundle := []tasks.Task{
		harness.ReadStatus,
		cluster.Delete,
	}

	return tasks.Run(bundle, nil)
}
