# Non-resource assertions

Under certain circumstances, it makes sense to evaluate assertions that do not depend on resources.
For example, when asserting the number of nodes in a cluster is equal to a known value.

The test below uses the `x_k8s_list` function to query the list of nodes in the cluster.
It uses the results to compare the number of nodes found with a known number (`1` in this case).

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          (x_k8s_list($client, 'v1', 'Node')):
            (length(items)): 1
```
