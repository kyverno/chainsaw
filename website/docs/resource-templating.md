# Use resource templating

Chainsaw simplifies dynamic resource configuration with native resource templating support.

Sometimes things we need to create resources or assertions are only known at runtime.

In the past, users have created all sorts of hacks, using tools like `envsubst` for dynamic substitution of env-variables. Those workarounds usually lack flexibility and introduce new problems like hiding the real resource to Chainsaw, preventing it to cleanup resources properly.

## Templating and Bindings

The templating engine in Chainsaw is based on the concept of **bindings**.

You can think of bindings as a side context where you can store and retrieve data based on keys. A resource template can read from the data from the side context to hydrate a concrete resource from the template.

Chainsaw offers some built-in bindings you can use. You can also create your own bindings and use outputs to pass information from one operation to the next.

!!! note
    Under the hood, Chainsaw uses the jmespath language, bindings are implemented using [lexical scoping](https://github.com/jmespath-community/jmespath.spec/blob/main/jep-011a-lexical-scope.md).

## Built-in bindings

The `$namespace` is a built-in binding provided by Chainsaw, containing the name of the ephemeral namespace used to execute a test (by default Chainsaw will create an ephemeral namespace for each test).

In the template below, we are using the `$namespace` binding at two different places, effectively injecting the ephemeral namespace name in the resource name and the `data.foo` field:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: template
spec:
  steps:
  - assert:
      resource:
        # apiVersion, kind, name, namespace and labels are considered for templating
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ($namespace)
        # other fields are not (they are part of the assertion tree)
        data:
          foo: ($namespace)
```

## Custom bindings

Built-in bindings allow templates to know about the context they are running in. On top of that, you can also create your own bindings, combining other bindings together, calling jmespath functions and so on.

In the template below we create bindings at different levels in a test and combine them by calling a jmespath function to configure an environment variable that will be available is script:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: script-env
spec:
  # bindings can be declared at the test level
  bindings:
  - name: chainsaw
    value: chainsaw
  steps:
    # bindings can also be declared at the step level
  - bindings:
    - name: hello
      value: hello
    try:
    - script:
        # bindings can also be declared at the operation level
        bindings:
        - name: awesome
          value: awesome
        env:
        - name: GREETINGS
          # multiple bindings can used to invoke a jmespath function
          value: (join(' ', [$hello, $chainsaw, 'is', $awesome]))
        content: echo $GREETINGS
```

!!! tip
    Bindings are immutable. This means two bindings can have the same name without overwriting each other.
    When a binding is registered it potentially hides other bindings with the same name. When this binding goes out of scope, previously registered bindings become visible again.

## Next step

Combining bindings and templates with operation outputs allows even more flexibility to [pass arbitrary data from one operation to another](./operation-outputs.md).
