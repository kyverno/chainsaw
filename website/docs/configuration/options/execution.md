# Execution options

Execution options determine how tests are run.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `failFast` | `false` | FailFast determines whether the test should stop upon encountering the first failure. |
| `parallel` | `auto` | The maximum number of tests to run at once. |
| `repeatCount` | `1` | RepeatCount indicates how many times the tests should be executed. |
| `forceTerminationGracePeriod` | | ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments. |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  execution:
    failFast: true
    parallel: 8
    repeatCount: 5
    forceTerminationGracePeriod: 5s
```

### With flags

TODO
```bash
chainsaw test                                   \
  --test-file chainsaw-test                     \
  --full-name                                   \
  --include-test-regex 'chainsaw/.*'            \
  --exclude-test-regex 'chainsaw/exclude-.*'
```
