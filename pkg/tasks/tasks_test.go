package tasks_test

import (
	"fmt"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2e-harness/pkg/tasks"
	"github.com/spf13/afero"
)

var (
	files   = []string{"/task1", "/task2"}
	taskErr = func(s *harness.Status) (*harness.Status, error) {
		return nil, fmt.Errorf("my-error")
	}
)

func getTaskFunc(filename string, fs afero.Fs) tasks.Task {
	return func(s *harness.Status) (*harness.Status, error) {
		if err := afero.WriteFile(fs, filename, []byte("test!"), 0644); err != nil {
			return nil, err
		}
		return s, nil
	}
}

func TestRunNoError(t *testing.T) {
	fs := new(afero.MemMapFs)

	bundle := []tasks.Task{getTaskFunc(files[0], fs), getTaskFunc(files[1], fs)}

	status := &harness.Status{}
	err := tasks.Run(bundle, status)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	for _, file := range files {
		e, err := afero.Exists(fs, file)
		if err != nil {
			t.Errorf("unexpected error %s", err)
		}
		if !e {
			t.Errorf("expected file %s to exists", file)
		}
	}
}

func TestRunError(t *testing.T) {
	fs := new(afero.MemMapFs)

	var bundle []tasks.Task
	bundle = append(bundle, getTaskFunc(files[0], fs))
	bundle = append(bundle, taskErr)
	bundle = append(bundle, getTaskFunc(files[1], fs))

	status := &harness.Status{}
	err := tasks.Run(bundle, status)
	if err == nil {
		t.Error("expected error didn't happen")
	}
	if err.Error() != "my-error" {
		t.Error("expected error didn't happen")
	}

	e, err := afero.Exists(fs, files[0])
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if !e {
		t.Errorf("expected file %s to exists", files[0])
	}

}

func TestRunIgnoreError(t *testing.T) {
	fs := new(afero.MemMapFs)

	var bundle []tasks.Task
	bundle = append(bundle, getTaskFunc(files[0], fs))
	bundle = append(bundle, taskErr)
	bundle = append(bundle, getTaskFunc(files[1], fs))

	status := &harness.Status{}
	tasks.RunIgnoreError(bundle, status)

	for _, file := range files {
		e, err := afero.Exists(fs, file)
		if err != nil {
			t.Errorf("unexpected error %s", err)
		}
		if !e {
			t.Errorf("expected file %s to exists", file)
		}
	}
}
