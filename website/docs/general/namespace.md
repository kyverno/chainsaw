# Test namespace

By default, Chainsaw will create an ephemeral namespace with a random name for each test, unless a specific namespace name is provided at the global or test level.

## Selection

### Global

One way to control the namespace used to run tests is to specify the name in the Chainsaw configuration [Namespace options](../configuration/options/namespace.md).

If a namespace name is specified at the configuration level Chainsaw will use it to run the tests (unless an individual test overrides the namespace name).

### Per test

If a namespace is specified in a test spec, Chainsaw will use it to run the test regardless of whether a namespace name was configured at the global level.

### Random

If no namespace name was specified at the global or test level, Chainsaw will create a random one for the lifetime of the test.

## Cleanup

As with any other resource, Chainsaw will clean up the namespace only if the namespace was created by Chainsaw.

If the namespace already exists when the test starts, Chainsaw will use it to run the test but won't delete it after the test terminates.

## Template

A namespace [template](./templating.md) can be provided at the global or test level.

This is useful if you want to make something specific with the namespace Chainsaw creates (add labels, add annotations, etc...).

!!! tip
    A namespace template specified at the test level takes precedence over the namespace template specified at the global level.

## Namespace injection

Because the name of the namespace is only known at runtime, depending on the resource being manipulated, Chainsaw will eventually inject the namespace name, except if:

- the resource already has a namespace specified
- the resource is a clustered resource

### Example

The resource below is a namespaced one and has no namespace specified.
Chainsaw will automatically inject the namespace name in it:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
  # there is no namespace configured and the resource
  # is a namespaced one.
  # Chainsaw will automatically inject the test namespace
data:
  foo: bar
```
