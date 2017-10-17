package harness

import (
	"github.com/giantswarm/micrologger"
)

type Harness struct {
	logger micrologger.Logger
}

func New(logger micrologger.Logger) *Harness {
	return &Harness{logger: logger}
}

func (h *Harness) Run() error {
	// create cluster if needed

	// run project setup

	// install sonobuoy infra

	// execute plugin

	// run project teardown

	return nil
}
