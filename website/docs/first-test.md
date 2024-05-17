# Create a Test

To create a chainsaw test all you need to do is to create one (or more) YAML file.

The recommended approach is to create one folder per test, with a `chainsaw-test.yaml` file containing one (or more) test definitions.
The test definition can reference other files in the same folder or anywhere else on the file system as needed.

!!! tip
    While chainsaw supports [other syntaxes](./tests/index.md), we strongly recommend the explicit approach.

## What is a test?

To put it simply, a test can be represented as an ordered sequence of test steps.

In turn, a test step can be represented as an ordered sequence of operations.

## Let's write our first test

For this quick start, we will create a (very simple) `Test` with one step and two operations:

1. Create a `ConfigMap` from a manifest
1. Verify the `ConfigMap` was created and contains the expected data

Follow the instructions below to create the folder and files defining our first test.

### Create a test folder

```bash
# create test folder
mkdir chainsaw-quick-start

# enter test folder
cd chainsaw-quick-start
```

### Create a `ConfigMap` manifest

```
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

### Create a test manifest

By default, Chainsaw will look for a file named `chainsaw-test.yaml` in every folder.

```
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

## Next step

Now we have created our first test, you can continue to the next section to [execute it](./run-tests.md).
