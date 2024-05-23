# Templating

Chainsaw simplifies dynamic configuration with native templating support.

Sometimes things we need to create resources or assertions are only known at runtime.

In the past, users have created all sorts of hacks using tools like `envsubst` for dynamic substitution of env-variables.
Those workarounds usually lack flexibility and introduce new problems like hiding the real resources from Chainsaw, preventing it from cleaning resources properly.

## Syntax

!!! tip
    Resource templating is heavily based on [bindings](./bindings/index.md) and uses [JMESPath](https://jmespath.site/) language.

### With bindings

In the template below, we are using the `$namespace` binding at two different places, effectively injecting the ephemeral namespace name in the `name` and the `data.foo` fields:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - assert:
      resource:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ($namespace)
        data:
          foo: ($namespace)
```

### With JMESPath

In the template below, we are using the JMESPath [join](https://jmespath.org/proposals/functions.html#join) function to create a unique resource name:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - apply:
      resource:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: (join('-', [$namespace, 'cm']))
        data:
          foo: bar
```
