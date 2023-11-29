# Delete

The `delete` operation allows you to specify resources that should be deleted from the Kubernetes cluster before a particular test step is executed.

!!! tip "Reference documentation"
    The full structure of the `Delete` is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Delete).

## Usage in `Test`

Below is an example of using `delete` in a `Test` resource.

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
        - delete:
            ref:
              apiVersion: v1
              kind: Pod
              namespace: default
              name: my-test-pod
        # ...
    ```

## Usage in `TestStep`

Below is an example of using `delete` in a `TestStep` resource.

!!! example

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: TestStep
    metadata:
      name: example
    spec:
      try:
      # ...
      - delete:
          ref:
            apiVersion: v1
            kind: Pod
            namespace: default
            name: my-test-pod
      # ...
    ```

## Operation check

Below is an example of using an [operation check](./check.md#delete).

!!! example "With check"

    ```yaml
    # ...
    - delete:
        ref:
          apiVersion: v1
          kind: Pod
          namespace: default
          name: my-test-pod
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
    - delete:
        ref:
          apiVersion: v1
          kind: Pod
          namespace: default
          name: my-test-pod
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
    # ...
    ```
