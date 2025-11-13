# Use bindings

You can think of bindings as a side context where you can store and retrieve data based on keys.

This is particularly useful when some data is only known at runtime.
For example, to pass data from one operation to another, to implement resource templating, to fetch data from an external system, etc.

Chainsaw offers some built-in bindings you can directly use in your tests but you can also create your own bindings if needed.

## Inheritance

Bindings can be configured at the test, step or operation level.

All bindings configured at a given level are automatically inherited in child levels.

!!! info "JMESPath"
    Chainsaw uses the [JMESPath](https://jmespath.site/) language, and bindings are implemented using [lexical scoping](https://github.com/jmespath-community/jmespath.spec/blob/main/jep-011a-lexical-scope.md).

## Immutability

Bindings are immutable. This means two bindings can have the same name without overwriting each other.

When a binding is registered it potentially hides other bindings with the same name.

When this binding goes out of scope, previously registered bindings with the same name become visible again.

## Built-in bindings

The `$namespace` binding is a good example of a **built-in binding** provided by Chainsaw.
It contains the name of the ephemeral namespace used to execute a test (by default Chainsaw will create an ephemeral namespace for each test).

In the operation below, we are assigning the value of the `$namespace` binding to an environment variable, and `echo` it in a script:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        env:
          # assign the value of the `$namespace` binding
          # to the environment variable `FOO`
        - name: FOO
          value: ($namespace)
        content: echo $FOO
```

## External values

The `$values` binding contains values provided via `--values`, `--set`, or `--set-string` flags.

## Custom bindings

On top of built-in bindings, you can also create your own ones, combine bindings together, call JMESPath functions using bindings as arguments, etc.

In the test below we create **custom bindings** at different levels in the test, combine them by calling the `join` function, assign the result to an environment variable, and `echo` it in a script:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
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
          # combined bindings together using the `join` functions and
          # assign the result to the GREETINGS environment variable
        - name: GREETINGS
          value: (join(' ', [$hello, $chainsaw, 'is', $awesome]))
        content: echo $GREETINGS
```

## Next step

Let's see how bindings can be useful with [resource templating](./resource-templating.md).
