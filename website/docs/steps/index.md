# What is a test step

A test step is made of three main components used to dermine the actions Chainsaw will perform when executing the step.

1. The [try](./try.md) statement *(required)*
1. The [catch](./catch.md) statement *(optional)*
1. The [finally](./finally.md) statement *(optional)*

!!! tip "Reference documentation"
    The full structure of the `TestStep` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-TestStep).

## Test step lifecycle

!!! info "Test step lifecycle"

    1. The step starts executing operations in the `try` statement
    1. If an operation fails in the `try` statement
        1. If a `catch` statement is present, **all operations** and collectors are executed
    1. If a `finally` statement is present, **all operations** and collectors are executed

## Example

The test step below highlights the basic structure of test step containing all `try`, `catch` and `finally` statements.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      # try to apply a couple of resources
      try:
      - apply:
          file: resources.yaml
      # in case of error, collect pod logs
      catch:
      - podLogs: {}
      # in all cases, collect events in the test namespace
      finally:
      - events: {}
    ```
