# Quick start

To use Chainsaw you will need a [Kubernetes](https://kybernetes.io) cluster, Chainsaw won't create one for you.

In this Quick start we will use [kind](https://kind.sigs.k8s.io) but feel free to use the tool of your choice.

!!! warning "Not a cluster management tool"

    We consider this is not the responsibility of Chainsaw to manage clusters.

    There are plenty of solutions to create and manage local clusters that will do that better than Chainsaw.

## Create a [kind](https://kind.sigs.k8s.io) cluster

Please refer to the [kind install docs](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) to install it locally.

Once [kind](https://kind.sigs.k8s.io) creating a local cluster is as simple as running:

```bash
# define kind image
export KIND_IMAGE="kindest/node:v1.28.0"

# create cluster
kind create cluster --image $KIND_IMAGE
```

## Writing tests

A Chainsaw test is [made of YAML files in a folder](./tests/what-is-a-test.md).

YAML files can either contain raw manifests with a special file naming convention to identify the step operations.
This is useful to create test quickly but doesn't allow great flexibility.

Another option is to have a `chainsaw-test.yaml` file containing a `Test` resource. While more verbose, this offers full flexibility over the test and test steps configuration.

### Quick start

For this quick start, we will create a `Test` with one step and two operations:

1. Create a `ConfigMap` from a manifest
1. Verify the `ConfigMap` was created and contains the expected data

### Create the test folder

```bash
# create test folder
mkdir chainsaw-quick-start

# enter test folder
cd chainsaw-quick-start
```

### Create a `ConfigMap`

```bash
# create a ConfigMap
cat > configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
EOF
```

### Create the test

By default, Chainsaw will look for a file named `chainsaw-test.yaml` in every folder.

```bash
# create test file
cat > chainsaw-test.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: quick-start
spec:
  steps:
  - try:
    # first operation: create the config map
    - apply:
        # file is relative to the test folder
        file: configmap.yaml
    # second operation: verify the config map exists and contains the expected data
    - assert:
        # file is relative to the test folder
        file: configmap.yaml
EOF
```

## Run Chainsaw

We finished writing our first test, now we can run Chainsaw to verify the test runs as expected:

```bash
$ chainsaw test

Loading default configuration...
- Using test file: chainsaw-test.yaml
- TestDirs [.]
- SkipDelete false
- FailFast false
- ReportFormat ''
- ReportName 'chainsaw-report'
- Namespace ''
- FullName false
- IncludeTestRegex ''
- ExcludeTestRegex ''
- ApplyTimeout 5s
- AssertTimeout 30s
- CleanupTimeout 30s
- DeleteTimeout 15s
- ErrorTimeout 30s
- ExecTimeout 5s
Loading tests...
- quick-start (.)
Running tests...
=== RUN   chainsaw
=== PAUSE chainsaw
=== CONT  chainsaw
=== RUN   chainsaw/quick-start
=== PAUSE chainsaw/quick-start
=== CONT  chainsaw/quick-start
    | 10:30:13 | quick-start | @setup   | CREATE    | OK    | v1/Namespace @ chainsaw-strong-troll
    | 10:30:13 | quick-start | step-1   | TRY       | RUN   |
    | 10:30:13 | quick-start | step-1   | APPLY     | RUN   | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | step-1   | CREATE    | OK    | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | step-1   | APPLY     | DONE  | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | step-1   | ASSERT    | RUN   | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | step-1   | ASSERT    | DONE  | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | step-1   | TRY       | DONE  |
    | 10:30:13 | quick-start | @cleanup | DELETE    | RUN   | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | @cleanup | DELETE    | OK    | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | @cleanup | DELETE    | DONE  | v1/ConfigMap @ chainsaw-strong-troll/chainsaw-quick-start
    | 10:30:13 | quick-start | @cleanup | DELETE    | RUN   | v1/Namespace @ chainsaw-strong-troll
    | 10:30:13 | quick-start | @cleanup | DELETE    | OK    | v1/Namespace @ chainsaw-strong-troll
    | 10:30:18 | quick-start | @cleanup | DELETE    | DONE  | v1/Namespace @ chainsaw-strong-troll
--- PASS: chainsaw (0.00s)
    --- PASS: chainsaw/quick-start (5.26s)
PASS
Tests Summary...
- Passed  tests 1
- Failed  tests 0
- Skipped tests 0
Done.
```
