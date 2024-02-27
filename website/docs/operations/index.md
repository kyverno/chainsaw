# Operations

While tests are made of test steps, test steps can be considered made of operations.

Every operation in a test step runs sequentially.

!!! warning "Only one action per operation"

    Every operation consists of a single action. While it is syntactically possible to create an operation with multiple actions, Chainsaw will verify and reject tests if operations containing multiple actions are found.

    The reasoning behind this intentional choice is that it becomes harder to understand in which order actions will be executed in case an operation consists of multiple actions. For this reason, operations consisting of multiple actions are disallowed.

### Common fields

All operations share some configuration fields.

!!! tip "Reference documentation"
    The full structure of the `Operation` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Operation).

#### ContinueOnError

Determines whether a test step should continue or not in case the operation is not successful.

!!! note ""
    Even if the test continues executing, it will still be reported as failed.

## Available operations

- [Apply](./apply.md)
- [Assert](./assert.md)
- [Command](./command.md)
- [Create](./create.md)
- [Delete](./delete.md)
- [Error](./error.md)
- [Patch](./patch.md)
- [Script](./script.md)
- [Sleep](./sleep.md)
- [Wait](./wait.md)

## Non-resource assertions

It is possible to evaluate assertions that do not depend on resources.

See [Non-resource assertions](./non-resource-assert.md) for details.

## Operation checks

Some operations support checking the operation execution result against specific expectations.

See [Operation checks](./check.md) for use case details and supported operations.

## Resource templating

Chainsaw can apply transformations to the resources before they are processed by the operation.

See [Resource templating](./templating.md) for use case details and supported operations.
