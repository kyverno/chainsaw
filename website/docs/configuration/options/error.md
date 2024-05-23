# Error options

Error options contains the global error configuration.

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
    Error options can not be configured with flags.
