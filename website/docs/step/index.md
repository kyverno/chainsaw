# What is a test step

A test step is made of three main components used to determine the actions Chainsaw will perform when executing the step.

1. The [try](./try.md) statement *(required)*
1. The [catch](./catch.md) statement *(optional)*
1. The [finally](./finally.md) statement *(optional)*

!!! tip "Reference documentation"
    The full structure of the `TestStep` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-TestStep).

## Test step lifecycle

!!! info "Test step lifecycle"

    1. The step starts executing operations in the `try` statement
    1. If an operation fails in the `try` statement
        1. If a `catch` statement is present, **all operations** and collectors are executed
    1. If a `finally` statement is present, **all operations** and collectors are executed

## Example

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: path/to/apply.yaml
        - assert:
            file: path/to/assert.yaml
        catch: []
        finally: []
    ```
