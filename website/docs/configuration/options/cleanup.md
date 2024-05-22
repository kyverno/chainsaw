# Cleanup options

Cleanup contains the cleanup configuration.

## Supported elements

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `skipDelete` | `bool` |  |  | <p>If set, do not delete the resources after running a test.</p> |
| `delayBeforeCleanup` | [`meta/v1.Duration`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration) |  |  | <p>DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.</p> |

## Configuration

### With file

```yaml
```

### With flags

```bash
```



## Delay before cleanup

At the end of each test, Chainsaw will delete the resources it created during the test.

When testing operators, it can be useful to wait a little bit before starting the cleanup process to make sure the operator/controller has the necessary time to update its internal state.

For this reason, Chainsaw provides the `delayBeforeCleanup` configuration option and the corresponding `--delay-before-cleanup` flag.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  # ...
  delayBeforeCleanup: 5s
  # ...
```

## Flag

```bash
chainsaw test --delay-before-cleanup 5s ...
```

TODO: skipDelete
