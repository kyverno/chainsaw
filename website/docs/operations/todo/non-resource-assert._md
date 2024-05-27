# Non-resource assertions

Under certain circumstances, it makes sense to evaluate assertions that do not depend on resources.
For example, when asserting the number of nodes in a cluster is equal to a known value.

## Usage examples

Below is an example of using `assert` in a `Test` resource.

!!! example "Using a file"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: non-resource-assertion
    spec:
      steps:
      - try:
        - assert:
            resource:
              (x_k8s_list($client, 'v1', 'Node')):
                (length(items)): 1
        - error:
            resource:
              (x_k8s_list($client, 'v1', 'Node')):
                (length(items)): 2
    ```
