module github.com/giantswarm/e2e-harness

go 1.13

require (
	github.com/giantswarm/apiextensions v0.2.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.16.6
	k8s.io/apiextensions-apiserver v0.16.6
	k8s.io/apimachinery v0.16.6
	k8s.io/client-go v0.16.6
	k8s.io/kube-aggregator v0.0.0
)

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191114103820-f023614fb9ea
