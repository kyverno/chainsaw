# Namespace options

Namespace options contain the configuration used to allocate a namespace for each test.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `name` | | Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec. |
| `template` | | Template defines a template to create the test namespace. |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  namespace:
    name: foo
    template:
      metadata:
        annotations:
          from-config-file: hello
```

### With flags

!!! note
    The `template` element can not be configured with flags.

```bash
chainsaw test --namespace foo
```
