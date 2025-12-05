# Patch

The `patch` operation defines resources that should be modified in a Kubernetes cluster.

If the resource to be modified does not exist in the cluster, the step will fail.

## Configuration

The full structure of `Patch` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Patch).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark: |
| [Outputs](../general/outputs.md) support           | :white_check_mark: |
| [Templating](../general/templating.md) support     | :white_check_mark: |
| [Operation checks](../general/checks.md) support   | :white_check_mark: |

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # use a specific file
        file: my-configmap.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # use glob pattern
        file: "configs/*.yaml"
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
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
    - patch:
        # specify resource inline
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: chainsaw-quick-start
          data:
            foo: bar
```

### Operation check

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        file: my-configmap.yaml
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: ConfigMap
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
```
