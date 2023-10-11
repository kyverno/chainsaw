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

### Create the first step

```bash
# create step file (note that the file name starts with `01`, it indicates the sequence order to run steps)
cat > 01-create-resource.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
spec:
  apply:
  - configmap.yaml
EOF
```

### Create the second step

```bash
# create step file (note that the file name starts with `02`, it indicates the sequence order to run steps)
cat > 02-check-resource-exists.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
spec:
  assert:
  - configmap.yaml
EOF
```

## Run Chainsaw

We finished writing our first test, now we can run Chainsaw to verify the test runs as expected:

```bash
chainsaw ...
```

TODO
