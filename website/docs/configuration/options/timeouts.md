# Timeouts

Timeouts in Chainsaw are specified per type of operation.
This is required because the timeout varies greatly depending on the nature of an operation.

For example, applying a manifest in a cluster is expected to be reasonably fast, while validating a resource can be a long operation.

## Supported timeouts

| Element | Default | Description |
|---|---|---|
| apply | `5s` | Used when Chainsaw applies manifests in a cluster |
| assert | `30s` | Used when Chainsaw validates resources in a cluster |
| cleanup | `30s` | Used when Chainsaw removes resources created for a test |
| delete | `15s` | Used when Chainsaw deletes resources from a cluster |
| error | `30s` | Used when Chainsaw validates resources in a cluster |
| exec | `5s` | Used when Chainsaw executes arbitrary commands or scripts |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  timeouts:
    apply: 45s
    assert: 20s
    cleanup: 45s
    delete: 25s
    error: 10s
    exec: 45s
```

### With flags

```bash
chainsaw test               \
  --apply-timeout 45s       \
  --assert-timeout 45s      \
  --cleanup-timeout 45s     \
  --delete-timeout 45s      \
  --error-timeout 45s       \
  --exec-timeout 45s
```
