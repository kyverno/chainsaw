# Catch

A `catch` statement is also a sequence of [operations](../operations/index.md).

Operations contained in a `catch` statement will be executed only if the step failed when executing the operations in the step's [try](./try.md) statement.

!!! tip
    All operations of a `catch` statement will be executed regardless of the success or failure of each of them.

## Operations

A `catch` statement supports only the following [operations](../operations/index.md):

- [Command](../operations/command.md)
- [Delete](../operations/delete.md)
- [Describe](../operations/helpers/describe.md)
- [Events](../operations/helpers/events.md)
- [Get](../operations/helpers/get.md)
- [Pod logs](../operations/helpers/logs.md)
- [Script](../operations/script.md)
- [Sleep](../operations/sleep.md)
- [Wait](../operations/helpers/wait.md)

## Inheritance

Under certain circumstances, it can be useful to configure catch blocks at a higher level than the step grain. At the test or configuration level.

This allows for declaring common catch statements we want to execute when an error occurs.
Those catch blocks are combined to produce the final catch block in the following order:

1. catch statements from the configuration level are executed first (if any)
1. catch statements from the test level are executed next (if any)
1. catch statements from the step level are executed last (if any)
