# Concurrency control

By default, Chainsaw will run tests in parallel.

The number of concurrent tests can be configured globally using a configuration file or with the `--parallel` flag.

Alternatively, the concurrent nature of a test can specified at the test level:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  # concurrency can be specified per test (`true` or `false`)
  # default value is `true`
  concurrent: true
  # ...
```

All non-concurrent tests are executed first, followed by the concurrent tests running in parallel.
