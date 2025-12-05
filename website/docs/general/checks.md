# Operation checks

Considering an operation's success or failure is not always as simple as checking an error code.

- Sometimes an operation can fail but the failure is what you expected, hence the operation should be reported as successful.
- Sometimes an operation can succeed but the result is not what you expected, in this case, the operation should be reported as a failure.

To support those kinds of use cases, some operations support additional checks to evaluate the operation result against an [assertion tree](https://kyverno.github.io/kyverno-json/latest/intro/).

## Input model

Different operations have a different model passed through the assertion tree.

Please consult the [Built-in bindings](../reference/builtins.md) reference documentation to learn what is available depending on the operation.

## Expect vs Check

While a simple check is enough to determine the result of a single operation, we needed a more advanced construct to cover `apply`, `create`, `delete`, `patch` and `update` operations. Those operations can operate on files containing [multiple resources](./references.md) and every resource can lead to a different result and expectation.

### Check

The example below uses a simple check. The operation is expected to fail (`($error != null): true`):

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        content: |
          exit 1
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
```

### Expect

To support more granular checks we use the `expect` field that contains an array of [Expectations](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Expectation).

Every expectation is made of an optional `match` combined with a `check` statement.

This way it is possible to control the scope of a `check`:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - create:
        file: resources.yaml
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: ConfigMap
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
```

In the test above, only config maps are expected to fail. If the `resources.yaml` file contains other type of resources they are supposed to be created without error (if an error happens for a non config map resource, the operation will be considered a failure).
