# Inline resources

When an operation needs to reference a resource, it can do so using a file path or directly specify the resource inline using the `resource` field.

The test below is equivalent to our [first test](../quick-start/first-test.md):

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
    - assert:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
```
