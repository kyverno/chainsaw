# Script

The `script` operation provides a means to run a script during the test step.

## Configuration

The full structure of the `Script` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark: |
| [Outputs](../general/outputs.md) support           | :white_check_mark: |
| [Templating](../general/templating.md) support     | :x:                |
| [Operation checks](../general/checks.md) support   | :white_check_mark: |

### KUBECONFIG

- Unless `--no-cluster` is specified, Chainsaw always executes commands in the context of a temporary `KUBECONFIG`, built from the configured target cluster.
- This specific `KUBECONFIG` has a single cluster, auth info and context configured (all named `chainsaw`).

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        content: |
          echo "hello chainsaw"
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
    - script:
        content: |
          echo "hello chainsaw"
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
```
