# Execution options

Execution options determine how tests are run by Chainsaw.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `failFast` | `false` | FailFast determines whether the test should stop upon encountering the first failure. |
| `parallel` | `auto` | The maximum number of tests to run at once. |
| `repeatCount` | `1` | RepeatCount indicates how many times the tests should be executed. |
| `forceTerminationGracePeriod` | | ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments. |

### Termination grace period

Some Kubernetes resources can take time before being terminated. For example, deleting a pod can take time if the underlying container doesn't quit quickly enough.

Chainsaw can override the grace period for the following resource kinds:

- Pod
- Deployment
- StatefulSet
- DaemonSet
- Job
- CronJob

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
    repeatCount: 2
    forceTerminationGracePeriod: 5s
```

### With flags

```bash
chainsaw test                                   \
  --fail-fast                                   \
  --parallel 8                                  \
  --repeat-count 2                              \
  --force-termination-grace-period 5s
```
