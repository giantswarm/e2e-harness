package tasks

import "github.com/giantswarm/e2e-harness/pkg/harness"

// Task represent a generic step in a pipeline
type Task func(*harness.Status) (*harness.Status, error)

func Run(tasks []Task, status *harness.Status) error {
	var err error
	for _, task := range tasks {
		status, err = task(status)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunIgnoreError(tasks []Task, status *harness.Status) {
	for _, task := range tasks {
		status, _ = task(status)
	}
}
