# Script

The `script` operation provides a means to run a script during the test step.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Script` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).
    - This operation supports [bindings](../tests/common/bindings.md).
    - This operation supports [outputs](../tests/common/outputs.md).

## Usage examples

Below is an example of using `script` in a `Test` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - script:
            content: |
              echo "hello chainsaw"
        # ...
    ```

## Operation check

Below is an example of using an [operation check](./check.md#script).

!!! example "With check"

    ```yaml
    # ...
    - script:
        content: |
          echo "hello chainsaw"
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
    # ...
    ```
