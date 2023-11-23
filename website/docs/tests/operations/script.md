# Script

The `script` operation provides a means to run a script during the test step.

!!! tip "Reference documentation"
    The full structure of the `Script` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

## Usage in `Test`

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

## Usage in `TestStep`

Below is an example of using `script` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
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
        (error != null): true
    # ...
    ```
