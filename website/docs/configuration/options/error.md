# Error options

Error options contain the global error configuration used by Chainsaw.

## Supported elements

| Field | Default | Description |
|---|---|---|
| `catch` | | Catch defines what the tests steps will execute when an error happens. This will be combined with catch handlers defined at the test and step levels. |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  error:
    catch:
    - events: {}
    - describe:
        resource: crds
```

### With flags

!!! note
    Error options can't be configured with flags.
