# Proxy

Proxy runs a proxy request against a configured Kubernetes cluster, targeting pods or services.

## Configuration

The full structure of `Proxy` is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Proxy).

### Features

| Supported features                                    |                    |
|-------------------------------------------------------|:------------------:|
| [Bindings](../../general/bindings.md) support         | :x:                |
| [Outputs](../../general/outputs.md) support           | :x:                |
| [Templating](../../general/templating.md) support     | :x:                |
| [Operation checks](../../general/checks.md) support   | :x:                |

### API version

The target resource API version to send the request to.

### Kind

The target resource kind to send the request to.

### Namespace

The target resource namespace to send the request to.

### Name

The target resource name to send the request to.

### Port

The target port to send the request to.

### Path

The target path to send the request to.

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  skip: true
  steps:
  - try:
    - proxy:
        # proxy request to the `kyverno-svc-metrics` service in the `kyverno` namespace
        apiVersion: v1
        kind: Service
        namespace: kyverno
        name: kyverno-svc-metrics
        # proxy request to the `metrics-port` port of the service
        port: metrics-port
        # send request to the `/metrics` path
        path: /metrics
        outputs:
          # decode received metrics and create an output with the results
        - name: metrics
          value: (x_metrics_decode($stdout))
```
