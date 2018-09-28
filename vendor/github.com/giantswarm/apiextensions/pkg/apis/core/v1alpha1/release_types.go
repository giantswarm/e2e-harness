package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewReleaseCRD returns a new custom resource definition for Release. This
// might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: releases.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: Release
//         plural: releases
//         singular: release
//       subresources:
//         status: {}
//
func NewReleaseCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "releases.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "Release",
				Plural:   "releases",
				Singular: "release",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ReleaseSpec   `json:"spec"`
	Status            ReleaseStatus `json:"status"`
}

type ReleaseSpec struct {
	Operator      ReleaseSpecOperator      `json:"operator" yaml:"operator"`
	VersionBundle ReleaseSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type ReleaseSpecOperator struct {
	Name    string `json:"cluster" yaml:"name"`
	Version string `json:"node" yaml:"version"`
}

type ReleaseSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type ReleaseStatus struct {
	Conditions []ReleaseStatusCondition `json:"conditions" yaml:"conditions"`
}

// ReleaseStatusCondition expresses a condition in which a node may is.
type ReleaseStatusCondition struct {
	// Status may be True, False or Unknown.
	Status string `json:"status" yaml:"status"`
	// Type may be Pending, Ready, Draining, Drained.
	Type string `json:"type" yaml:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Release `json:"items"`
}
