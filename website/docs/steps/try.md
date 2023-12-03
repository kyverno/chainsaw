# Try

A `try` statement is a sequence of [operations](../operations/index.md) executed in the same order they are declared.
If an operation fails the entire step is considered failed.

The `try` statement is at the heart of a test step, it represents what the step is supposed to be about.

[catch](./catch.md) and [finally](./finally.md) statements should be viewed as complementary to the `try` statement.

!!! tip "Continue on error"

    By default a test step stops executing when an operation fails and the following operations are not executed.

    This behaviour can be changed using the `continueOnError` field, if `continueOnError` is set to `true` the step will still be considered failed but execution will continue with the next operations.

## Operations

A `try` statement supports all [operations](../operations/index.md):

- [Apply](../operations/apply.md)
- [Assert](../operations/assert.md)
- [Command](../operations/command.md)
- [Create](../operations/create.md)
- [Delete](../operations/delete.md)
- [Error](../operations/error.md)
- [Script](../operations/script.md)
- [Sleep](../operations/sleep.md)

## Example

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # apply a couple of resources in the cluster
      - apply:
          file: resources.yaml
      # evaluate a couple of assertions
      - assert:
          file: asserts.yaml
      # take a small break
      - sleep:
          duration: 15s
      # execute a custom script
      - script:
          content: echo "Hello Chainsaw"
      # ...
    ```
