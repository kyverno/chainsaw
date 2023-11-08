# Create

The `create` operation lets you define resources that should be created in the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

!!! warning

    If the resource to be created already exists in the cluster, the step will fail.

The full structure of the `Create` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Create).


## Usage in `Test`

Below is an example of using `create` in a `Test` resource.

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
        - create:
            file: ../resources/configmap.yaml
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `create` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - create:
          file: ../resources/configmap.yaml
      # ...
    ```
