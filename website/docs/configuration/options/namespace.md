# Namespace options

Namespace contains info about the namespace used for testing.

## Supported elements

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` |  |  | <p>Name defines the namespace to use for tests. If not specified, every test will execute in a random ephemeral namespace unless the namespace is overridden in a the test spec.</p> |
| `template` | `policy/v1alpha1.Any` |  |  | <p>Template defines a template to create the test namespace.</p> |

## Configuration

### With file

```yaml
```

### With flags

```bash
```
