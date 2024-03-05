# Finally

A `finally` statement is similar to a [catch](./catch.md) statement but will always execute after the [try](./try.md) and eventual [catch](./catch.md) statements finished executing regardless of the success or failure of the test step.

!!! tip
    All operations and collectors of a `finally` statement will be executed regardless of the success or failure of each of them.

## Operations

A `finally` statement supports only the following [operations](../operations/index.md):

- [Command](../operations/command.md)
- [Script](../operations/script.md)
- [Sleep](../operations/sleep.md)
- [Delete](../operations/delete.md)
- [Wait](../operations/wait.md)

## Collectors

A `finally` statement supports all [collectors](../collectors/index.md):

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
        catch: []
        finally:
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
        - wait:
            # ...
    ```
