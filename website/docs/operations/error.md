# Error

The `error` operation lets you define a set of expected errors for a test step.
If any of these errors occur during the test, they are treated as expected outcomes.
However, if an error that's not on this list occurs, it will be treated as a test failure.

## Configuration

The full structure of `Error` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Error).

### Features

| Supported features                                 |                           |
|----------------------------------------------------|:-------------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark:        |
| [Outputs](../general/outputs.md) support           | :x:                       |
| [Templating](../general/templating.md) support     | :x: \| :white_check_mark: |
| [Operation checks](../general/checks.md) support   | :x:                       |

### Templating

When working with `assert` and `error` operations, the content is already an assertion tree and therefore mostly represents a logical operation. An exception to this rule is for fields participating in the resource selection process.

For this reason, only elements used for looking up the resources from the cluster will be considered for templating. That is, `apiVersion`, `kind`, `name`, `namespace` and `labels`.

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
