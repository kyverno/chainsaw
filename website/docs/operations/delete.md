# Delete

The `delete` operation defines resources that should be deleted from a Kubernetes cluster.

!!! warning
    The propagation policy is forced to `Background` because some types default to `Orphan` (this is the case for unmanaged jobs for example) and we don't want to let dangling pods run in the cluster after cleanup.

## Configuration

The full structure of `Delete` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Delete).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark: |
| [Outputs](../general/outputs.md) support           | :x:                |
| [Templating](../general/templating.md) support     | :white_check_mark: |
| [Operation checks](../general/checks.md) support   | :white_check_mark: |

### Propagation policy

The deletion propagation policy can be overriden for the specific operation using the `deletionPropagationPolicy` field.

When the `deletionPropagationPolicy` is not specified it is [inherited](../general/inheritance.md) from the configuration, test, or test step.

Look at the [configuration section](../configuration/options/deletion.md#propagation) for details about supported values and how this affects execution.

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - delete:
        # specify reference inline
        ref:
          apiVersion: v1
          kind: Pod
          namespace: default
          name: my-test-pod
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - delete:
        # use a file
        file: my-configmap.yaml
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
    - delete:
        ref:
          apiVersion: v1
          kind: Pod
          namespace: default
          name: my-test-pod
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: Pod
            metadata:
              namespace: default
              name: my-test-pod
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
```
