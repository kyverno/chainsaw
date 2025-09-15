# Cleanup

A `cleanup` statement is similar to a [finally](./finally.md) statement but will execute after the test finishes executing, while `finally` executes after the step finishes executing.

!!! tip
    All operations of a `cleanup` statement will be executed regardless of the success or failure of each of them.

## Operations

A `cleanup` statement supports only the following [operations](../operations/index.md):

- [Command](../operations/command.md)
- [Delete](../operations/delete.md)
- [Describe](../operations/helpers/describe.md)
- [Events](../operations/helpers/events.md)
- [Get](../operations/helpers/get.md)
- [Pod logs](../operations/helpers/logs.md)
- [Script](../operations/script.md)
- [Sleep](../operations/sleep.md)
- [Wait](../operations/helpers/wait.md)
