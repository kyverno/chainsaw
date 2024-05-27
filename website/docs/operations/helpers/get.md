# Get

Display one or many resources.

## Configuration

The full structure of the `Get` resource is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Get).

!!! warning "Deprecated syntax"
    You can specify the `resource` directly instead of using `apiVersion` and `kind`.
    
    **This is a deprecated syntax though and will be removed in a future version.**

## Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

## Test namespace

When used with a namespaced resource, Chainsaw will default the scope to the ephemeral test namespace.

## All namespaces

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
    # get all pods in the test namespace
    - get:
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
    - get:
        apiVersion: v1
        kind: Pod
        # get pods that have a name starting with the provided `my-pod`
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
    - get:
        apiVersion: v1
        kind: Pod
        # get pods in the namespace `foo`
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
    - get:
        apiVersion: v1
        kind: Pod
        # get pods using a label selector query
        selector: app=my-app
```

### Format

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - get:
        apiVersion: v1
        kind: Pod
        format: json
```
