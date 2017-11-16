package results

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/giantswarm/e2e-harness/pkg/runner"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
)

const (
	DefaultResultsFilename    = "results.xml"
	DefaultResultsPath        = "/workdir/plugins/e2e/" + DefaultResultsFilename
	DefaultRemoteResultsPath  = "/tmp/results"
	DefaultTarResultsFilename = "e2e.tar.gz"
)

type TestSuite struct {
	Tests     int    `xml:"tests,attr,omitempty"`
	Failures  int    `xml:"failures,attr,omitempty"`
	Errors    int    `xml:"errors,attr,omitempty"`
	Time      string `xml:"time,attr,omitempty"`
	TestCases []TestCase
}

type TestCase struct {
	Name    string       `xml:"name,attr"`
	Error   *TestFailure `xml:"error,omitempty"`
	Failure *TestFailure `xml:"failure,omitempty"`
}

// TestFailure contains data related to a failed test.
type TestFailure struct {
	Value   string `xml:",innerxml"`
	Type    string `xml:"type,attr,omitempty"`
	Message string `xml:"message,attr,omitempty"`
}

type Results struct {
	logger micrologger.Logger
	fs     afero.Fs
	runner runner.Runner
}

func New(logger micrologger.Logger, fs afero.Fs, r runner.Runner) *Results {
	return &Results{
		logger: logger,
		fs:     fs,
		runner: r,
	}
}

func (r *Results) Read(path string) (*TestSuite, error) {
	cmd := "cat " + path
	b := new(bytes.Buffer)
	if err := r.runner.RunPortForward(b, cmd); err != nil {
		return nil, microerror.Mask(err)
	}

	ts := &TestSuite{}

	if err := xml.Unmarshal(b.Bytes(), ts); err != nil {
		return nil, microerror.Mask(err)
	}
	return ts, nil
}

func Write(fs afero.Fs, results *TestSuite) error {
	if err := fs.MkdirAll(path.Dir(DefaultRemoteResultsPath), os.ModePerm); err != nil {
		return microerror.Mask(err)
	}

	content, err := xml.Marshal(results)
	if err != nil {
		return microerror.Mask(err)
	}

	resultsFilename := filepath.Join(DefaultRemoteResultsPath, DefaultResultsFilename)
	err = afero.WriteFile(fs, resultsFilename, []byte(content), 0644)
	if err != nil {
		return microerror.Mask(err)
	}

	doneFilename := filepath.Join(DefaultRemoteResultsPath, "done")
	err = afero.WriteFile(fs, doneFilename, []byte(resultsFilename), 0644)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Results) Unpack() error {
	cmds := []string{
		"rm -f /workdir/results/*.tar.gz",
		"kubectl cp heptio-sonobuoy/sonobuoy:/tmp/sonobuoy /workdir/results --namespace=heptio-sonobuoy",
		"tar xzf /workdir/results/*.tar.gz",
	}
	for _, cmd := range cmds {
		if err := r.runner.RunPortForward(ioutil.Discard, cmd); err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

// Interpret is a Task that knows how to grab test results and extract
// the execution outcome from them.
func (r *Results) Interpret() error {
	ts, err := r.Read(DefaultResultsPath)
	if err != nil {
		return microerror.Mask(err)
	}

	if ts.Failures == 0 && ts.Errors == 0 {
		r.logger.Log("info", "sonobuoy plugin succeeded")
		return nil
	}
	return fmt.Errorf("sonobuoy plugin failed")
}
