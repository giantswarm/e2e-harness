package harness

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Status represents the current state of the
// test harness. It is the token passed amongst
// tasks.
type Status struct {
	RemoteCluster bool   `yaml:"remoteCluster"`
	GitCommit     string `yaml:"gitCommit"`
}

// Init is a Task for initializing the harness.
func Init(status Status) (Status, error) {
	baseDir, err := BaseDir()
	if err != nil {
		return status, err
	}
	err = os.MkdirAll(filepath.Join(baseDir, "workdir"), os.ModePerm)

	return status, err
}

// WriteStatus is a Task that persists the input status to a file.
func WriteStatus(status Status) (Status, error) {
	dir, err := BaseDir()
	if err != nil {
		return status, err
	}

	content, err := yaml.Marshal(&status)
	if err != nil {
		return status, err
	}

	if err := ioutil.WriteFile(
		filepath.Join(dir, "status.yaml"),
		[]byte(content),
		0644); err != nil {
		return status, err
	}

	return status, nil
}

// ReadStatus is a Task that populates a Status struct with the data read
// from a default file location.
func ReadStatus(Status) (Status, error) {
	dir, err := BaseDir()
	if err != nil {
		return Status{}, err
	}

	content, err := ioutil.ReadFile(filepath.Join(dir, "status.yaml"))
	if err != nil {
		return Status{}, err
	}

	s := &Status{}

	if err := yaml.Unmarshal(content, s); err != nil {
		return Status{}, err
	}

	return *s, nil
}

func BaseDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".e2e-harness"), nil
}
