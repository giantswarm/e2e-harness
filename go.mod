module github.com/giantswarm/e2e-harness

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions v0.4.4
	github.com/giantswarm/apprclient v0.2.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/helmclient v1.0.2
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	k8s.io/kube-aggregator v0.17.3
)
