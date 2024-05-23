# Apply

The `apply` operation lets you define resources that should be applied to the Kubernetes cluster during the test step.

These can be configurations, deployments, services, or any other Kubernetes resource.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Apply` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Apply).
    - This operation supports [bindings](../general/bindings.md).
    - This operation supports [outputs](../general/outputs.md).

## Usage examples

Below is an example of using `apply` in a `Test` resource.

!!! example "Using a specific file"

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
            file: my-configmap.yaml
        # ...
    ```

!!! example "Using file path expressions"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example-multi
    spec:
      steps:
      - try:
        # ...
        - apply:
            file: "configs/*.yaml"
        # ...
    ```

!!! example "Using an URL"

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
            file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/step/configmap.yaml
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

### Operation check

Below is an example of using an [operation check](./check.md#apply).

!!! example "With check"

    ```yaml
    # ...
    - apply:
        file: my-configmap.yaml
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: ConfigMap
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
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
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: ConfigMap
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
    # ...
    ```
