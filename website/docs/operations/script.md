# Script

The `script` operation provides a means to run a script during the test step.

## Configuration

The full structure of the `Script` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

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
