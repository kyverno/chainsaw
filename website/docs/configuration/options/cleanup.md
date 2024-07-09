# Cleanup options

Cleanup options contain the configuration used by Chainsaw for cleaning up resources.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `skipDelete` | `false` | If set, do not delete the resources after running a test. |
| `delayBeforeCleanup` | | DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts. |

### Delay before cleanup

At the end of each test, Chainsaw will delete the resources it created during the test.

When testing operators, it can be useful to wait a little bit before starting the cleanup process to make sure the operator/controller has the necessary time to update its internal state.

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  cleanup:
    skipDelete: true
    delayBeforeCleanup: 5s
```

### With flags

```bash
chainsaw test                   \
  --skip-delete                 \
  --cleanup-delay 5s
```
