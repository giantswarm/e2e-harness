package framework

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/giantswarm/microerror"
)

// HelmCmd executes a helm command.
func HelmCmd(cmd string) error {
	err := runCmd("helm " + cmd)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func runCmd(cmdStr string) error {
	log.Printf("Running command %q\n", cmdStr)
	cmdEnv := os.ExpandEnv(cmdStr)
	fields := strings.Fields(cmdEnv)
	cmd := exec.Command(fields[0], fields[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
