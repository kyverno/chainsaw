# Update

The `update` operation lets you define resources that should be updated in the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

If the resource to be updated doesn't exist in the cluster, the step will fail.

## Configuration

The full structure of the `Update` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Update).

!!! tip
    - This operation supports [bindings](../general/bindings.md).
    - This operation supports [outputs](../general/outputs.md).

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - update:
        # use a specific file
        file: my-configmap.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example-multi
spec:
  steps:
  - try:
    - update:
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
    - update:
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
    - update:
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
    - update:
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
