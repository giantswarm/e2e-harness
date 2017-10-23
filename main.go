package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/e2e-harness/cmd"
)

var (
	gitCommit = "latest"
	name      = "e2e-harness"
)

func main() {
	cmd.SetGitCommit(gitCommit)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
