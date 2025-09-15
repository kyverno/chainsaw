# Assert

The `assert` operation allows you to specify conditions that should hold true for a successful test.

For example, after applying resources, you might want to ensure that a particular pod is running or a service is accessible.

## Configuration

The full structure of `Assert` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Assert).

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
