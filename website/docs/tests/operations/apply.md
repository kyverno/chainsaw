# Apply

The `apply` operation lets you define resources that should be applied to the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

The full structure of the `Apply` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Apply).

## Usage in `Test`

Below is an example of using `apply` in a `Test` resource.

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
        - apply:
            file: my-pod.yaml
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `apply` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - apply:
          file: my-pod.yaml
      # ...
    ```
