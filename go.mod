module github.com/giantswarm/e2e-harness

go 1.13

require (
	github.com/giantswarm/apiextensions v0.2.0
	github.com/giantswarm/apprclient v0.2.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/helmclient v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.16.6
	k8s.io/apiextensions-apiserver v0.16.6
	k8s.io/apimachinery v0.16.6
	k8s.io/apiserver v0.16.6 // indirect
	k8s.io/client-go v0.18.0
	k8s.io/helm v2.16.4+incompatible
	k8s.io/kube-aggregator v0.0.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191114105449-027877536833
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5-beta.1
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191114103151-9ca1dc586682
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191114110141-0a35778df828
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191114112024-4bbba8331835
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191114111741-81bb9acf592d
	k8s.io/code-generator => k8s.io/code-generator v0.16.5-beta.1
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191114102325-35a9586014f7
	k8s.io/cri-api => k8s.io/cri-api v0.16.5-beta.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191114112310-0da609c4ca2d
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191114103820-f023614fb9ea
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191114111510-6d1ed697a64b
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191114110717-50a77e50d7d9
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191114111229-2e90afcb56c7
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191114113550-6123e1c827f7
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191114110954-d67a8e7e2200
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191114112655-db9be3e678bb
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191114105837-a4a2842dc51b
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191114104439-68caf20693ac
)