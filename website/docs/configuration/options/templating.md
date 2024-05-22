# Templating options

Templating contains the templating configuration.

## Supported elements

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `enabled` | `bool` |  |  | <p>Enabled determines whether resources should be considered for templating.</p> |

## Configuration

### With file

```yaml
```

### With flags

```bash
```











!!! warning "Experimental status"
    This is an **experimental feature**, and implementation could change slightly in the next versions.

!!! info
    Templating was disabled by default in `v0.1.*` but is now enabled by default since `v0.2.1`.

Chainsaw can apply transformations to the resources before they are processed by the operation.

This is useful when a resource needs some runtime configuration.

Templating must be enabled at the configuration, test, step, or operation level for the templating process to kick in.
Alternatively, templating can be enabled using the `--template` flag when invoking chainsaw from the command line.

!!! note
    Unlike assertion trees, templating can only be specified in leave nodes of the YAML tree.

## Example

The test below will create, assert, and delete a `ConfigMap` with a dynamic `name` configured at runtime using the `$namespace` binding.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: values
spec:
  steps:
  - try:
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: ($namespace)
    - assert:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: ($namespace)
    - delete:
        ref:
          apiVersion: v1
          kind: ConfigMap
          name: ($namespace)
```
