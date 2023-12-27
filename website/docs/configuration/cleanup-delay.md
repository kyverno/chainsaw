# Delay before cleanup

At the end of each test, Chainsaw will delete the resources it created during the test.

When testing operators, it can be useful to wait a little bit before starting the cleanup process to make sure the operator/controller has the necessary time to update the internal state.

For this reason, Chainsaw provides the `delayBeforeCleanup` configuration option and the corresponding `--delay-before-cleanup` flag.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: custom-config
spec:
  # ...
  delayBeforeCleanup: 5s
  # ...
```

## Flag

```bash
$ chainsaw test --delay-before-cleanup 5s ...
```
