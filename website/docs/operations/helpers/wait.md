# Wait

Wait for a specific condition on one or many resources.

## Configuration

The full structure of the `Wait` resource is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Wait).

### Features

| Supported features                                    |                    |
|-------------------------------------------------------|:------------------:|
| [Bindings](../../general/bindings.md) support         | :x:                |
| [Outputs](../../general/outputs.md) support           | :x:                |
| [Templating](../../general/templating.md) support     | :x:                |
| [Operation checks](../../general/checks.md) support   | :x:                |

### Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

### All resources

If you don't specify a `name` or a `selector`, the `wait` operation will consider `all` resources.

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
  - try:
    # wait all pods are ready in the test namespace
    - wait:
        apiVersion: v1
        kind: Pod
        timeout: 1m
        for:
          condition:
            name: Ready
            value: 'true'
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        # wait a specific pod is ready in the test namespace
        name: my-pod
        timeout: 1m
        for:
          condition:
            name: Ready
            value: 'true'
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        # wait all pods are ready in the namespace `foo`
        namespace: foo
        timeout: 1m
        for:
          condition:
            name: Ready
            value: 'true'
```

### Label selector

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        # match pods using a label selector query
        selector: app=foo
        timeout: 1m
        for:
          condition:
            name: Ready
            value: 'true'
```

### Deletion

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        timeout: 1m
        for:
          # wait for deletion
          deletion: {}
```

### JSON Path

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        timeout: 1m
        for:
          # arbitrary JSON path
          jsonpath:
            path: '{.status.phase}'
            value: Running
```

### Format

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - wait:
        apiVersion: v1
        kind: Pod
        format: json
```
