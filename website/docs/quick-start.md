# Quick start

To use Chainsaw you will need a [Kubernetes](https://kybernetes.io) cluster, Chainsaw won't create one for you.

We consider this is not the responsibility of Chainsaw to manage clusters.
There are plenty of solutions to create and manage local clusters that will do that better than Chainsaw.

In this Quick start we will use [kind](https://kind.sigs.k8s.io) but feel free to use the tool of your choice.

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

A Chainsaw test is made of YAML files in a folder.

Every file contains a `TestStep` and Chainsaw will run every step sequentially.

For this Quick start, we will create a two step test:

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

```bash
# create test file
cat > chainsaw-test.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: quick-start
spec:
  steps:
  # first step applies the config map
  - spec:
      apply:
      # file is relative to the test folder
      - file: configmap.yaml
  # second step verifies the config map exists and contains the expected data
  - spec:
      assert:
      # file is relative to the test folder
      - file: configmap.yaml
EOF
```

## Run Chainsaw

We finished writing our first test, now we can run Chainsaw to verify the test runs as expected:

```bash
chainsaw test --test-dir .

Loading default configuration...
- Timeout 30s
- TestDirs [.]
- SkipDelete false
- FailFast false
- Parallel 8
- ReportFormat ''
- ReportName 'chainsaw-report'
- Namespace ''
- FullName false
- IncludeTestRegex ''
- ExcludeTestRegex ''
Loading tests...
- quick-start (.)
Running tests...
=== RUN   quick-start
=== PAUSE quick-start
=== CONT  quick-start
    logger.go:17: 12:31:16 | quick-start | CREATE[v1/Namespace] chainsaw-feasible-thrush
    logger.go:17: 12:31:16 | quick-start/step-1 | APPLY[v1/ConfigMap] chainsaw-feasible-thrush/chainsaw-quick-start
    logger.go:17: 12:31:16 | quick-start/step-1 | CREATE[v1/ConfigMap] chainsaw-feasible-thrush/chainsaw-quick-start
    logger.go:17: 12:31:16 | quick-start/step-2 | ASSERT[v1/ConfigMap] chainsaw-feasible-thrush/chainsaw-quick-start
    logger.go:17: 12:31:16 | quick-start/step-1 | DELETE[v1/ConfigMap] chainsaw-feasible-thrush/chainsaw-quick-start
    logger.go:17: 12:31:16 | quick-start | DELETE[v1/Namespace] chainsaw-feasible-thrush
--- PASS: quick-start (5.22s)
PASS
Tests Summary...
- Passed  tests 1
- Failed  tests 0
- Skipped tests 0
Done.
```
