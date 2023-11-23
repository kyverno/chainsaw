# Apply

The `apply` operation lets you define resources that should be applied to the Kubernetes cluster during the test step.

These can be configurations, deployments, services, or any other Kubernetes resource.

!!! tip "Reference documentation"
    The full structure of the `Apply` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Apply).

## Usage in `Test`

Below is an example of using `apply` in a `Test` resource.

!!! example "Using a file"

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

!!! example "Using an inline resource"

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
            resource:
              apiVersion: v1
              kind: ConfigMap
              metadata:
                name: chainsaw-quick-start
              data:
                foo: bar
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `apply` in a `TestStep` resource.

!!! example "Using a file"

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

!!! example "Using an inline resource"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - apply:
          resource:
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: chainsaw-quick-start
            data:
              foo: bar
      # ...
    ```

## Operation check

Below is an example of using an [operation check](./check.md#apply).

!!! example "With check"

    ```yaml
    # ...
    - apply:
        file: my-pod.yaml
      check:
        # an error is expected, this will:
        # - succeed if the operation failed
        # - fail if the operation succeeded
        (error != null): true
    # ...
    ```

!!! example "With check"

    ```yaml
    # ...
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: chainsaw-quick-start
          data:
            foo: bar
      check:
        # an error is expected, this will:
        # - succeed if the operation failed
        # - fail if the operation succeeded
        (error != null): true
    # ...
    ```
