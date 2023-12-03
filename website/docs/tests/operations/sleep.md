# Sleep

The `sleep` operation provides a means to sleep for a configured duration.

!!! tip "Reference documentation"
    The full structure of the `Sleep` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Sleep).

## Usage in `Test`

Below is an example of using `sleep` in a `Test` resource.

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
        - sleep:
            duration: 30s
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `sleep` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - sleep:
          duration: 30s
      # ...
    ```
