# Command

The `command` operation provides a means to execute a specific command during the test step.

The full structure of the `Command` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Command).

## Usage in `Test`

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

## Usage in `TestStep`

Below is an example of using `command` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - command:
          entrypoint: echo
          args:
          - hello chainsaw
      # ...
    ```
