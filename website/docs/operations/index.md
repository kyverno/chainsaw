# Operations

Chainsaw supports the following operations:

- [Apply](./apply.md)
- [Assert](./assert.md)
- [Command](./command.md)
- [Create](./create.md)
- [Delete](./delete.md)
- [Error](./error.md)
- [Patch](./patch.md)
- [Script](./script.md)
- [Sleep](./sleep.md)
- [Update](./update.md)

## Helpers

Chainsaw also supports [kubectl helpers](./helpers/index.md).

## Properties

### Action unicity

Every operation must consist of a single action.

While it is syntactically possible to create an operation with multiple actions, Chainsaw will verify and reject tests if operations containing multiple actions are found.

The reasoning behind this intentional choice is that it becomes harder to understand in which order actions will be executed when an operation consists of multiple actions. For this reason, operations consisting of multiple actions are not allowed.

## Common fields

### Continue on error

The `continueOnError` field determines whether a test step should continue executing or not if the operation fails (in any case the test will be marked as failed).

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
      # in case of error the test will be marked as failed
      # but the step will not stop execution and will
      # continue executing the following operations
    - continueOnError: true
      apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
```

### Description

All operations support a `description` field that can be used document your tests.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - description: Waits a couple of seconds
      sleep:
        duration: 3s
```
