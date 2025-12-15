# Events

Display one or many events.

## Configuration

The full structure of `Events` is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Events).

### Features

| Supported features                                    |                    |
|-------------------------------------------------------|:------------------:|
| [Bindings](../../general/bindings.md) support         | :x:                |
| [Outputs](../../general/outputs.md) support           | :x:                |
| [Templating](../../general/templating.md) support     | :x:                |
| [Operation checks](../../general/checks.md) support   | :x:                |

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
    # get all events in the test namespace
    - events: {}
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    # get events in a specific namespace
    - events:
        namespace: foo
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    # get event by name
    - events:
        name: my-event
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
    - events:
        # get events using a label selector query
        selector: app=my-app
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try: ...
    catch:
    - events:
        # get events using a label selector query
        selector: app=my-app
        namespace: foo
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
    - events:
        format: json
```
