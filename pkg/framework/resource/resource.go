package resource

import (
	"fmt"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/apprclient"
	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/helm/pkg/helm"
)

type ResourceConfig struct {
	ApprClient *apprclient.Client
	HelmClient *helmclient.Client
	Logger     micrologger.Logger

	Namespace string
}

type Resource struct {
	apprClient *apprclient.Client
	helmClient *helmclient.Client
	logger     micrologger.Logger

	namespace string
}

func New(config ResourceConfig) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.ApprClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ApprClient must not be empty", config)
	}
	if config.HelmClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.Namespace == "" {
		config.Namespace = "giantswarm"
	}
	c := &Resource{
		apprClient: config.ApprClient,
		helmClient: config.HelmClient,
		logger:     config.Logger,

		namespace: config.Namespace,
	}

	return c, nil
}

func (r *Resource) InstallResource(name, values, version string, conditions ...func() error) error {
	chartValuesEnv := os.ExpandEnv(values)
	chartname := fmt.Sprintf("%s-chart", name)

	tarball, err := r.apprClient.PullChartTarballFromRelease(chartname, version)
	if err != nil {
		return microerror.Mask(err)
	}
	err = r.helmClient.InstallFromTarball(tarball, r.namespace, helm.ReleaseName(name), helm.ValueOverrides([]byte(chartValuesEnv)), helm.InstallWait(true))
	if err != nil {
		return microerror.Mask(err)
	}

	for _, c := range conditions {
		err = backoff.Retry(c, framework.NewExponentialBackoff(framework.ShortMaxWait, framework.ShortMaxInterval))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
