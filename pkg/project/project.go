package project

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/giantswarm/e2e-harness/pkg/docker"
	"github.com/giantswarm/e2e-harness/pkg/patterns"
	"github.com/giantswarm/micrologger"
	yaml "gopkg.in/yaml.v2"
)

type E2e struct {
	Version  string `yaml:"version"`
	Setup    []Step `yaml:"setup"`
	Teardown []Step `yaml:"teardown"`
}

type Step struct {
	Run     string   `yaml:"run"`
	WaitFor WaitStep `yaml:"waitFor"`
}

type WaitStep struct {
	Run     string        `yaml:"run"`
	Match   string        `yaml:"match"`
	Timeout time.Duration `yaml:"timeout"`
}

type Project struct {
	logger micrologger.Logger
	docker *docker.Docker
}

const (
	defaultTimeout = 120
)

var (
	initTiller = Step{
		Run: "helm init",
		WaitFor: WaitStep{
			Run:   "kubectl get pod -n kube-system",
			Match: `tiller-deploy.*1/1\s*Running`,
		},
	}
)

func New(logger micrologger.Logger, docker *docker.Docker) *Project {
	return &Project{
		logger: logger,
		docker: docker,
	}
}

func (p *Project) CommonSetupSteps() error {
	p.logger.Log("info", "executing common setup steps")
	steps := []Step{initTiller}
	for _, step := range steps {
		if err := p.runStep(step); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) SetupSteps() error {
	p.logger.Log("info", "executing setup steps")

	e2e, err := p.readProjectFile()
	if err != nil {
		return err
	}

	for _, step := range e2e.Setup {
		if err := p.runStep(step); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) TeardownSteps() error {
	p.logger.Log("info", "executing teardown steps")

	e2e, err := p.readProjectFile()
	if err != nil {
		return err
	}

	for _, step := range e2e.Teardown {
		if err := p.runStep(step); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) runStep(step Step) error {
	p.logger.Log("info", fmt.Sprintf("executing step with command %q", step.Run))
	if err := p.docker.RunPortForward(os.Stdout, step.Run); err != nil {
		return err
	}

	if err := p.wait(step.WaitFor); err != nil {
		return err
	}
	return nil
}

func (p *Project) wait(ws WaitStep) error {
	if ws.Timeout == 0 {
		ws.Timeout = defaultTimeout
	}

	timeout := time.After(ws.Timeout * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout looking for pattern %q with command %q", ws.Match, ws.Run)
		case <-tick:
			// pipe the output of the docker command to the input of FindMatch,
			// this way we can handle potentially long outputs witout having
			// to store them in a variable
			r, w := io.Pipe()
			// writing without a reader will deadlock so write in a goroutine
			go func() {
				defer w.Close()
				p.docker.RunPortForward(w, ws.Run)
			}()
			p.logger.Log("debug", "checking pattern "+ws.Match)
			ok, err := patterns.FindMatch(p.logger, r, ws.Match)
			if err != nil {
				return err
			} else if ok {
				p.logger.Log("debug", "match found")
				// Match found
				return nil
			}
			p.logger.Log("debug", "match not found, retrying")
		}
	}
}

func (p *Project) readProjectFile() (*E2e, error) {
	// read project file
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	projectFile := filepath.Join(dir, "e2e", "project.yaml")
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		return nil, err
	}

	content, err := ioutil.ReadFile(projectFile)
	if err != nil {
		return nil, err
	}

	e2e := &E2e{}

	if err := yaml.Unmarshal(content, e2e); err != nil {
		return nil, err
	}
	return e2e, nil
}
