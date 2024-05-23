# Command

The `command` operation provides a means to execute a specific command during the test step.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Command` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Command).
    - This operation supports [bindings](../tests/common/bindings.md).
    - This operation supports [outputs](../tests/common/outputs.md).

## Usage examples

Below is an example of using `command` in a `Test` resource.

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
        - command:
            entrypoint: echo
            args:
            - hello chainsaw
        # ...
    ```

> When defining shell command `args` in YAML format, it's crucial to consider potential differences in behavior, as Chainsaw may interpret them differently compared to regular shell or bash environments, due to quote removal.

## Operation check

Below is an example of using an [operation check](./check.md#command).

!!! example "With check"

    ```yaml
    # ...
    - command:
        entrypoint: echo
        args:
        - hello chainsaw
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
    # ...
    ```
