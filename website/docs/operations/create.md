# Create

The `create` operation lets you define resources that should be created in the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

!!! warning

    If the resource to be created already exists in the cluster, the step will fail.

!!! tip "Reference documentation"
    The full structure of the `Create` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Create).


## Usage in `Test`

Below is an example of using `create` in a `Test` resource.

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
        - create:
            file: my-configmap.yaml
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
        - create:
            file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/valid.yaml
        # ...

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
        - create:
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

Below is an example of using an [operation check](./check.md#create).

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
