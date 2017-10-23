package cluster

import (
	"os"

	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/harness"
)

// Create is a Task that creates a remote cluster.
func Create(status harness.Status) (harness.Status, error) {
	return clusterAction(status, "/usr/local/bin/shipyard", "-action=start")
}

// Delete is a Task that gets rid of a remote cluster.
func Delete(status harness.Status) (harness.Status, error) {
	return clusterAction(status, "/usr/local/bin/shipyard", "-action=stop")
}

func clusterAction(status harness.Status, args ...string) (harness.Status, error) {
	if !status.RemoteCluster {
		return status, nil
	}
	err := docker.Run(status.GitCommit, os.Stdout, args[0], args[1:]...)

	return status, err
}
