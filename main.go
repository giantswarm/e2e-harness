package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/e2e-harness/v2/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
