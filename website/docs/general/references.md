# References

Chainsaw tests often need to reference resources. Including references in tests can be done in multiple ways.

## Inline

One way to declare a resource is to do it directly inside the test definition:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # specify resource inline
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
```

This doesn't encourage file reuse but can be handy, especially when the resource definition is short or when the execution environment doesn't support file system access.

## File reference

Another option is to use the `file` field. The `file` can be a specific file, or multiple files declared using a glob pattern:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # use a specific file
        file: my-configmap.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # use glob pattern
        file: "configs/*.yaml"
```

## URL reference

A third option is to use a URL. Chainsaw uses https://github.com/hashicorp/go-getter, it will download the content from the remote service and load it in the operation resources:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # use an URL
        file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/step/configmap.yaml
```

## Cardinality

When using file-based references, it is important to note that the referenced file(s) can declare multiple resources. Internally, Chainsaw will duplicate the operation once per resource.

This is important to keep this in mind, especially when working with bindings and outputs. Bindings and outputs will be evaluated for every operation instance.
