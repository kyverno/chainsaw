# Describe

Show details of a specific resource or group of resources.

## Configuration

The full structure of `Describe` is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Describe).

### Features

| Supported features                                    |                    |
|-------------------------------------------------------|:------------------:|
| [Bindings](../../general/bindings.md) support         | :x:                |
| [Outputs](../../general/outputs.md) support           | :x:                |
| [Templating](../../general/templating.md) support     | :x:                |
| [Operation checks](../../general/checks.md) support   | :x:                |

### Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

### Test namespace

When used with a namespaced resource, Chainsaw will default the scope to the ephemeral test namespace.

### All namespaces

When used with a namespaced resource, it is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - describe:
        # describe all pods in the test namespace
        apiVersion: v1
        kind: Pod
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - describe:
        apiVersion: v1
        kind: Pod
        # describe pods that have a name starting with the provided `my-pod`
        name: my-pod
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - describe:
        apiVersion: v1
        kind: Pod
        # describe pods in the namespace `foo`
        namespace: foo
```

### Label selector

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - describe:
        apiVersion: v1
        kind: Pod
        # describe pods using a label selector query
        selector: app=my-app
```

### Show events

!!! tip
    By default, `showEvents`is `true`.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - describe:
        apiVersion: v1
        kind: Pod
        showEvents: false
```
