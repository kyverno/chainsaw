# Templating options

Templating options contain the templating configuration.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `enabled` | `true` | Enabled determines whether resources should be considered for templating. |

!!! tip
    Templating was disabled by default in `v0.1.*` but is now enabled by default since `v0.2.1`.

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  templating:
    enabled: false
```

### With flags

```bash
chainsaw test --template=false
```
