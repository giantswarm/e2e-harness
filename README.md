# e2e-harness

[![CircleCI](https://circleci.com/gh/giantswarm/e2e-harness.svg?style=shield)](https://circleci.com/gh/giantswarm/e2e-harness)

Harness for custom kubernetes e2e testing.

## Getting Project

Clone the git repository: https://github.com/giantswarm/e2e-harness.git

## Running e2e-harness

You can download a prebuilt binary from [here](https://github.com/giantswarm/e2e-harness/releases/) or,
with a golang environment set up, build from source from the root of the project:
```
go build .
```

## How does e2e-harness work

The goal of the project is making it easy to design and run e2e tests for kubernetes
components. We have great tools for local development, like minikube, but the tests
and results obtained using them are difficult to replicate on a CI environment.

e2e-harness aims to abstract all the differences between local and CI environments,
so that you can write the tests once and run them everywhere, making sure that if
things work locally they will work too elsewhere.

In order to achive that, e2e-harness has two operation modes: local and remote.
The setup and teardown actions differ on each mode, but the test themselves (and
the actions required to execute them) are the same.

Regarding the test execution, all the actions are run on a container, so that the
execution environment is always the same. We have put in place these binaries inside
the container:

* kubectl: k8s CLI client, allows us to run common setup actions or out of cluster
tests (see below).
* helm: at giantswarm most of the systems under tests are helm charts. We have
installed the registry plugin too.
* shipyard: it allows us to create and delete a remote minikube instance.

The container image is published in quay registry `quay.io/giantswarm/e2e-harness:latest`
and its Dockerfile can be found [here](https://github.com/giantswarm/e2e-harness/blob/master/Dockerfile).

## Requirements

The main requirement is having a recent docker version running on the host. Additionally,
for each operation mode:

* Local: [minikube](https://github.com/kubernetes/minikube) should be started with
RBAC enabled before running e2e-harness:

```
$ minikube start --extra-config=apiserver.Authorization.Mode=RBAC
```

* Remote: as stated above e2e-harness uses [shipyard](https://giathub.com/giantswarm/shipyard)
for setting up the remote cluster. shipyard currently only supports AWS as the
backend engine, so the common environment variables for granting access to AWS
are required too (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`).

## Project initialization

From the project root execute:

```
$ e2e-harness init
```

This will create a `e2e` directory, with the required files to start writing tests,
see below. This is how the e2e directory looks like this:

```
├── Dockerfile
├── main.go
├── project.yaml
└── tests
    ├── example.go
    └── runner.go
```

`project.yaml` defines how the end to end tests should be setup, which tests to
run and how to tear them down. It includes the following fields:
* `setup`: contains an array of steps to set the testing environment up.
* `outOfClusterTest`: set of steps to execute checks as an external client
* `inClusterTest`: configuration for the tests to be run as sonobuoy plugins, see
below.
* `teardown`: steps required to clean up the test environment.

The rest of the files are an example of in cluster tests. These are executed as
sonobuoy plugins. The example is written in go but you could use any langauge,
see [the sonobuoy examples section](https://github.com/heptio/sonobuoy/tree/master/examples). These are the files included:

* `Dockerfile`: Required, this is the image definition for the sonobuoy plugin that
runs the in-cluster tests, see below for details.
* `main.go`: Main file of the golang example, just calls the runner.
* `tests/runner.go`: Part of the golang example, helper file that defines methods
required for executing the tests.
* `tests/example.go`: Part of the golang example, this is where the actual test is
written. It must be registered in the runner during the `init` execution:
```
func init() {
  Add(TestExample)
}
```

The test functions need to fulfill this interface:

```
type Test func() (description string, err error)
```

The `description` used in the sonobuoy test output.

## e2e-harness lifecycle

An e2e test execution involves three stages:

* Setup: this is performed with the `setup` command.

```
$ e2e-harness setup --remote=[false|true]
```

It takes care of:
  - Initializing the project, creating an interchange directory that will
  be mounted on the test container for keeping state between command
  executions.
  - Prepare the connection to the test cluster: for remote executions this
  will involve the creation of the cluster too. In local executions, the
  connection settings to be able running minikube are made available in
  the test container interchange directory.
  - Run common set up steps: these will put the test cluster in a common
  initial state, basically installing tiller (helm's server side part) including
  the required resources to make it work with RBAC enabled.
  - Run specific setup steps: this are defined in the project file
  `e2e/project.yaml` under the `setup` key. They are common steps (see their
  description above) and can do things like installing the chart under test,
  setting up required external resources, etc.
* Run tests: invoking the `test` command:

```
$ e2e-harness test
```

First, the out of cluster tests are executed, if any, defined as
regular `step` entries under the `outOfClusterTests` key. Then the in cluster tests,
if they are enabled in `project.yaml`. You should enabled them only if they are
present, if not enabling them causes an error, see above for a descrition about
how to write this kind of tests.
* Teardown: the teardown phase is executing using the `teardown` command.

```
$ e2e-harness teardown
```

It includes:
  - Run specific teardown steps: this are defined in the project file
  `e2e/project.yaml` under the `teardown` key. They are common steps (see their
  description above) and can do things like removing the chart under test,
  tearing down required external resources, etc.
  - Run commomn tear down steps: these differ depending on the mode of operation,
  for remote, ephemeral clusters they are just deleted, for local clusters tiller
  and all the required RBAC setup is removed.

## Writing tests

e2e tests can be executed from agents either out of the test cluster or running
inside that cluster:

### Out of cluster tests

Currently the out of cluster tests are executed as `step` elements in the
`e2e/project.yaml` file, the elements in a `step` are:

* `run`: string with the command to execute, multiline allowed.
* `waitFor`: optional element to check if the `run` command was successful, it has
as nested items a `run` entry with the check command to be executed, a `pattern`
entry with a regular expression pattern to be checked against the previous command
output and a `timeout` with the number of miliseconds to keep executing the command
and checking the output. An example step could be:

```
- run: kubectl create namespace giantswarm
  waitFor:
    run: kubectl get namespace
    match: giantswarm\s*Active
    timeout: 2000
```

In this case, the first `run` entry tries to create a namespace and the `waitFor` entry
makes sure that the namespace is created and in Active state, with a deadline of 2 seconds.

All the `run` and `pattern` elements expand environment variables of the form `$ENV_VAR`
or `${ENV_VAR}`.

As all the steps in the project file, these elements are run from the test container,
so keep in mind that the binaries used must be available in it.

### In cluster tests

The in cluster tests execution are controlled by the `inClusterTests` entry in the project file:

```
inClusterTests:
  enabled: true
  env:
  - name: ENV_VAR_1
    value: VALUE_1
  - name: ENV_VAR_2
    value: VALUE_2
```

First of all, they must be enabled to be executed (set `enabled` field to `true`).
You can also pass the plugin environment variables to control their execution through
the `env` key.

The in cluster tests are executed as [sonobuoy plugins](https://github.com/heptio/sonobuoy/blob/master/docs/plugins.md).
e2e-harness takes care of installing all the required sonobuoy infrastructure,
running the plugins with a container created from the image defined by `e2e/Dockerfile`
grabbing the results and making them available in `.e2e-harness/workdir/plugin/e2e/results.xml`.
The example project deployed with `e2e-harness init` shows how this can be done
using golang, see the `Project initialilzation` section for details.

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/e2e-harness/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING.md](/giantswarm/e2e-harness/blob/master/CONTRIBUTING.md) for details on submitting patches, the contribution workflow as well as reporting bugs.

## License

E2e-Harness is under the Apache 2.0 license. See the [LICENSE](/giantswarm/e2e-harness/blob/master/LICENSE) file for details.
