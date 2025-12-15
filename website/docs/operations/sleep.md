# Sleep

The `sleep` operation provides a means to sleep for a configured duration.

## Configuration

The full structure of `Sleep` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Sleep).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :x:                |
| [Outputs](../general/outputs.md) support           | :x:                |
| [Templating](../general/templating.md) support     | :x:                |
| [Operation checks](../general/checks.md) support   | :x:                |

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
