package cmd

import (
	"fmt"
	"os"

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

const (
	// EnvVarK8sApiUrl is the process environment variable representing the
	// k8s api url for testing cluster.
	EnvVarK8sApiUrl = "K8S_API_URL"
	// EnvVarK8sCert is the process environment variable representing the
	// k8s kubeconfig cert value for testing cluster.
	EnvVarK8sCert = "K8S_CERT_ENCODED"
	// EnvVarK8sCert is the process environment variable representing the
	// k8s kubeconfig ca cert value for testing cluster.
	EnvVarK8sCertCA = "K8S_CERT_CA_ENCODED"
	// EnvVarK8sCert is the process environment variable representing the
	// k8s kubeconfig private key value for testing cluster.
	EnvVarK8sCertPrivate = "K8S_CERT_PRIVATE_ENCODED"
)

var (
	SetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "setup e2e tests",
		Run:   runSetup,
	}
)

var (
	setupCmdTestDir string
	name            string
	existingCluster bool
	remoteCluster   bool
)

func init() {
	RootCmd.AddCommand(SetupCmd)

	SetupCmd.Flags().StringVar(&setupCmdTestDir, "test-dir", project.DefaultDirectory, "Name of the directory containing executable tests.")
	SetupCmd.Flags().StringVar(&name, "name", "e2e-harness", "CI execution identifier")
	SetupCmd.Flags().BoolVar(&remoteCluster, "remote", true, "use remote cluster")
	SetupCmd.Flags().BoolVar(&existingCluster, "existing", false, "can be used with --remote=true to use already existing cluster")
}

func runSetup(cmd *cobra.Command, args []string) {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	err = runSetupError(cmd, args)
	if err != nil {
		logger.Log("level", "error", "message", "exiting with status 1 due to error", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}
}

func runSetupError(cmd *cobra.Command, args []string) error {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return microerror.Mask(err)
	}

	k8sApiUrl := os.Getenv(EnvVarK8sApiUrl)
	k8sCert := os.Getenv(EnvVarK8sCert)
	k8sCertCA := os.Getenv(EnvVarK8sCertCA)
	k8sCertPrivate := os.Getenv(EnvVarK8sCertPrivate)
	if existingCluster {
		if k8sApiUrl == "" {
			return microerror.Maskf(invalidConfigError, fmt.Sprintf("env var '%s' must not be empty", EnvVarK8sApiUrl))
		}

		if k8sCert == "" {
			return microerror.Maskf(invalidConfigError, fmt.Sprintf("env var '%s' must not be empty", EnvVarK8sCertCA))
		}

		if k8sCertCA == "" {
			return microerror.Maskf(invalidConfigError, fmt.Sprintf("env var '%s' must not be empty", EnvVarK8sCert))
		}

		if k8sCertPrivate == "" {
			return microerror.Maskf(invalidConfigError, fmt.Sprintf("env var '%s' must not be empty", EnvVarK8sCertPrivate))
		}
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
		ExistingCluster: existingCluster,
		RemoteCluster:   remoteCluster,
	}
	h := harness.New(logger, fs, hCfg)
	c := cluster.New(logger,
		fs,
		d,
		existingCluster,
		remoteCluster,
		k8sApiUrl,
		k8sCert,
		k8sCertCA,
		k8sCertPrivate)

	// tasks to run
	bundle := []tasks.Task{
		h.Init,
		h.WriteConfig,
		c.Create,
		p.CommonSetupSteps,
	}

	return microerror.Mask(tasks.Run(bundle))
}
