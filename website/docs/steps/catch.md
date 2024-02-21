# Catch

A `catch` statement is also a sequence of [operations](../operations/index.md) or [collectors](../collectors/index.md).

Operations and collectors contained in a `catch` statement will be executed only if the step failed when executing the operations in the step's [try](./try.md) statement.

!!!tip
    All operations and collectors of a `catch` statement will be executed regardless of the success or failure of each of them.

## Operations

A `catch` statement supports only the following [operations](../operations/index.md):

- [Command](../operations/command.md)
- [Script](../operations/script.md)
- [Sleep](../operations/sleep.md)
- [Delete](../operations/delete.md)

## Collectors

A `catch` statement supports all [collectors](../collectors/index.md):

- [Pod logs](../collectors/pod-logs.md)
- [Events](../collectors/events.md)
- [Get](../collectors/get.md)
- [Describe](../collectors/describe.md)

## Example

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try: []
        catch:
        - command:
            # ...
        - script:
            # ...
        - delete:
            # ...
        - events:
            # ...
        - podLogs:
            # ...
        - describe:
            # ...
        - get:
            # ...
        - sleep:
            # ...
        finally: []
    ```
