# Delete

The `delete` operation allows you to specify resources that should be deleted from the Kubernetes cluster before a particular test step is executed.

!!! info
    The propagation policy is forced to `Background` because some types default to `Orphan` (this is the case for unmanaged jobs for example) and we don't want to let dangling pods run in the cluster after cleanup.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Delete` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Delete).
    - This operation supports [bindings](../general/bindings.md).

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

Below is an example of using an [operation check](./todo/check.md#delete).

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
