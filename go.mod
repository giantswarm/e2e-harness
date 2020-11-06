module github.com/giantswarm/e2e-harness/v3

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions/v3 v3.7.0
	github.com/giantswarm/apprclient/v2 v2.0.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/helmclient/v3 v3.0.1
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.3
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/spf13/afero v1.4.1
	github.com/spf13/cobra v1.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.9
	k8s.io/apiextensions-apiserver v0.18.9
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
	k8s.io/kube-aggregator v0.18.5
)

replace sigs.k8s.io/cluster-api v0.3.10 => github.com/giantswarm/cluster-api v0.3.10-gs
