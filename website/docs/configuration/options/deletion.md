# Deletion options

Deletion options determine the configuration used by Chainsaw for deleting resources.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `propagation` | `Background` | Propagation decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. |

### Propagation

This element will affect [Kubernetes cascading deletion](https://kubernetes.io/docs/concepts/architecture/garbage-collection/#cascading-deletion).
Supported values are `Orphan`, `Background` and `Foreground`.

!!! tip
    Setting `Orphan` is probably never a good idea because it would leak resources in the test cluster. Chainsaw uses `Background` as its default value which is a reasonable choice.

    Note that `Foreground` can be useful to fail when the dependent resources fail to delete.

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  deletion:
    propagation: Foreground
```

### With flags

!!! note
    Deletion options can't be configured with flags.
