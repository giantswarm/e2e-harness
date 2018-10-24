package release

import (
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
)

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

func IsReleaseAlreadyExists(err error) bool {
	return helmclient.IsReleaseAlreadyExists(err)
}

func IsReleaseNotFound(err error) bool {
	return helmclient.IsReleaseNotFound(err)
}

var releaseStatusNotMatchingError = &microerror.Error{
	Kind: "releaseStatusNotMatchingError",
}

// IsReleaseStatusNotMatching asserts releaseStatusNotMatchingError
func IsReleaseStatusNotMatching(err error) bool {
	return microerror.Cause(err) == releaseStatusNotMatchingError
}

var releaseVersionNotMatchingError = &microerror.Error{
	Kind: "releaseVersionNotMatchingError",
}

// IsReleaseVersionNotMatching asserts releaseVersionNotMatchingError
func IsReleaseVersionNotMatching(err error) bool {
	return microerror.Cause(err) == releaseVersionNotMatchingError
}

func IsTarballNotFound(err error) bool {
	return helmclient.IsTarballNotFound(err)
}

func IsTillerNotFound(err error) bool {
	return helmclient.IsTillerNotFound(err)
}

var tooManyResultsError = &microerror.Error{
	Kind: "tooManyResultsError",
}

// IsTooManyResults asserts invalidConfigError.
func IsTooManyResults(err error) bool {
	return microerror.Cause(err) == tooManyResultsError
}

var unexpectedStatusPhaseError = &microerror.Error{
	Kind: "unexpectedStatusPhaseError",
}

// IsUnexpectedStatusPhase asserts notFoundError.
func IsUnexpectedStatusPhase(err error) bool {
	return microerror.Cause(err) == unexpectedStatusPhaseError
}
