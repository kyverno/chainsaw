# Delete

The `delete` operation allows you to specify resources that should be deleted from the Kubernetes cluster before a particular test step is executed.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Delete` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Delete).
    - This operation supports [bindings](../bindings/index.md).

## Usage examples

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
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: Pod
            metadata:
              namespace: default
              name: my-test-pod
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
    # ...
    ```
