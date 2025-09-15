# Pods logs

Print the logs for a container in a pod or specified resource.

## Configuration

The full structure of `PodLogs` is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-PodLogs).

### Features

| Supported features                                    |                    |
|-------------------------------------------------------|:------------------:|
| [Bindings](../../general/bindings.md) support         | :x:                |
| [Outputs](../../general/outputs.md) support           | :x:                |
| [Templating](../../general/templating.md) support     | :x:                |
| [Operation checks](../../general/checks.md) support   | :x:                |

### Test namespace

Chainsaw will default the scope to the ephemeral test namespace.

### All namespaces

It is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

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
    # all pods in the test namespace
    - podLogs: {}
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - podLogs:
        # pods that have a name starting with the provided `my-pod`
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
    - podLogs:
        # pods in the namespace `foo`
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
    - podLogs:
        # match pods using a label selector query
        selector: app=my-app
```

### Tail

!!! tip
    By default, `tail` will be `10` when a label selector is used, and `all` if a pod `name` is specified.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - podLogs:
        tail: 30
```

### Container

!!! tip
    By default logs from all containers will be fetched.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - podLogs:
        container: nginx
```
