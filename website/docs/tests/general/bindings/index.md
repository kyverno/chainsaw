# Bindings

You can think of bindings as a side context where you can store and retrieve data by name.

This is particularly useful when some data is only known at runtime.
For example, to pass data from one operation to another, to implement resource templating, to fetch data from an external system, etc.

Chainsaw offers some built-in bindings you can directly use in your tests, but you can also create your own bindings if needed.

## Syntax

!!! tip
    Browse the [reference documentation](../../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Binding) to see the syntax details and where bindings can be declared.

The test below illustrates bindings declaration at different levels:

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

## Inheritance

Bindings can be configured at the test, step or operation level.

All bindings configured at a given level are automatically [inherited](../inheritance.md) at lower levels.

!!! info "JMESPath"
    Chainsaw uses the [JMESPath](https://jmespath.site/) language, and bindings are implemented using [lexical scoping](https://github.com/jmespath-community/jmespath.spec/blob/main/jep-011a-lexical-scope.md).

## Immutability

Bindings are immutable. This means two bindings can have the same name without overwriting each other.

When a binding is registered it potentially hides other bindings with the same name.

When this binding goes out of scope, previously registered bindings with the same name become visible again.

## Templating

Both `name` and `value` of a binding can use [templating](../templating.md).

## Built-in bindings

Browse the [built-in bindings list](./builtins.md) to find available built-in bindings.
