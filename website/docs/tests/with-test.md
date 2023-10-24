# `Test`s based syntax

Chainsaw offers a flexible and precise way to define tests, especially tailored for Kubernetes-related operations. This guide assists in illustrating how to craft tests using Chainsaw's API.

## TestSpec: Defining the Core Test

`TestSpec` is the backbone of your test, allowing you to set the primary attributes of the test run:

1. `Timeout`: Determines how long the test should run before being marked as failed due to a timeout.

2. `Skip`: A simple flag to decide if a particular test should be ignored during the test run.

3. `Steps`: An ordered collection of actions or verifications (test steps) to be executed during the test.

## Writing tests

### 1. Setting Up the ConfigMap

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
  name: test-name
spec:
  timeout: "10m"  # Optional: Define the test's timeout duration.
  skip: false     # Optional: Decide if this test should be skipped.
  steps:
  # The first step applies the config map.
  - apply:
      - file: configmap.yaml  # Reference the ConfigMap we created.
  
  # The next step ensures the config map's existence and verifies its content.
  - assert:
      - file: configmap.yaml  # Assertion against the ConfigMap's expected state.
  
  # Additional steps can be added as per the requirements.

EOF
```
