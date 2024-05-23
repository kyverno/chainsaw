# Sleep

The `sleep` operation provides a means to sleep for a configured duration.

## Configuration

!!! tip "Reference documentation"
    The full structure of the `Sleep` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Sleep).

## Usage examples

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
