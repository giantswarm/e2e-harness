module github.com/giantswarm/e2e-harness

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions v0.4.17-0.20200723160042-89aed92d1080
	github.com/giantswarm/apprclient v0.2.1-0.20200724085653-63c7eb430dcf
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/helmclient v1.0.6-0.20200724131413-ea0311052b6e
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/spf13/afero v1.3.2
	github.com/spf13/cobra v1.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.5
	k8s.io/apiextensions-apiserver v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	k8s.io/kube-aggregator v0.18.5
)
