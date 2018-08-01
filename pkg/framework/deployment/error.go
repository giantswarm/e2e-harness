package deployment

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var wrongLabelsError = &microerror.Error{
	Kind: "wrongLabelsError",
}

// IsWrongLabels asserts wrongLabelsError.
func IsWrongLabels(err error) bool {
	return microerror.Cause(err) == wrongLabelsError
}

var wrongReplicasError = &microerror.Error{
	Kind: "wrongReplicasError",
}

// IsWrongReplicas asserts wrongReplicasError.
func IsWrongReplicas(err error) bool {
	return microerror.Cause(err) == wrongReplicasError
}
