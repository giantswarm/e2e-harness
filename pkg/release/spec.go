package release

import "context"

type ConditionFunc func() error

type ConditionSet interface {
	// SecretExists return a function waiting for the Secret to appear in the
	// Kubernetes API.
	SecretExists(ctx context.Context, namespace, name string) ConditionFunc
}
