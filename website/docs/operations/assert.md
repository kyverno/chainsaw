# Assert

The `assert` operation allows you to specify conditions that should hold true for a successful test.

For example, after applying resources, you might want to ensure that a particular pod is running or a service is accessible.

## Configuration

The full structure of the `Assert` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Assert).

!!! tip
    - This operation supports [bindings](../general/bindings.md).

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        # use a specific file
        file: ../resources/deployment-assert.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        # use glob pattern
        file: "../assertions/*.yaml"
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        # use an URL
        file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/valid.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        # specify resource inline
        resource:
          apiVersion: v1
          kind: Deployment
          metadata:
            name: foo
          spec:
            (replicas > 3): true
```
