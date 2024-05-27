# Sleep

The `sleep` operation provides a means to sleep for a configured duration.

## Configuration

The full structure of the `Sleep` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Sleep).

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - sleep:
        duration: 30s
```
