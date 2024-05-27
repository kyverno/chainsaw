# Error

The `error` operation lets you define a set of expected errors for a test step.
If any of these errors occur during the test, they are treated as expected outcomes.
However, if an error that's not on this list occurs, it will be treated as a test failure.

## Configuration

The full structure of the `Error` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Error).

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
    - error:
        # use a specific file
        file: ../resources/deployment-error.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - error:
        # use glob pattern
        file: "../errors/*.yaml"
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - error:
        # use an URL
        file: https://raw.githubusercontent.com/user/repo/branch/path/to/deployment-error.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - error:
        # specify resource inline
        resource:
          apiVersion: v1
          kind: Deployment
          metadata:
            name: foo
          spec:
            (replicas > 3): true
```
