# Operations

While tests are made of test steps, test steps can be considered made of operations.

Every operation in a test steps runs sequentially.

### Common fields

All operations share some configuration fields.

- **Timeout:** A timeout for the operation.
- **ContinueOnError:** Determines whether a test step should continue or not in case the operation was not successful.
  Even if the test continues executing, it will still be reported as failed.

The full structure of the `Operation` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Operation).

## Available operartions

- [Delete](./delete.md)
- [Apply](./apply.md)
- [Create](./create.md)
- [Assert](./assert.md)
- [Error](./error.md)
- [Command](./command.md)
- [Script](./script.md)
