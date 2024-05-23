# Outputs

Operation outputs can be useful for communicating and reusing computation results across operations.

Chainsaw evaluates outputs after an operation has finished executing. The results of output evaluations are registered in the bindings and are made available for the following operations.

## Syntax

!!! tip
    Browse the [reference documentation](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Output) to see the syntax details and where outputs can be declared.

### Basic

The test below illustrates output usage:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  bindings:
  - name: chainsaw
    value: chainsaw
  steps:
  - bindings:
    - name: hello
      value: hello
    try:
    - script:
        bindings:
        - name: awesome
          value: awesome
        env:
        - name: GREETINGS
          value: (join(' ', [$hello, $chainsaw, 'is', $awesome]))
        # output is used to register a new `$OUTPUT` binding
        outputs:
        - name: OUTPUT
          value: ($stdout)
        content: echo $GREETINGS
    - script:
        # output from the previous operation is used
        # to configure an evironment variable
        env:
        - name: INPUT
          value: ($OUTPUT)
        content: echo $INPUT
```

### Matching

An output supports an optional `match` field. The `match` is used to conditionally create the output binding.

The test below illustrates output with matching:

TODO

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  bindings:
  - name: chainsaw
    value: chainsaw
  steps:
  - bindings:
    - name: hello
      value: hello
    try:
    - script:
        bindings:
        - name: awesome
          value: awesome
        env:
        - name: GREETINGS
          value: (join(' ', [$hello, $chainsaw, 'is', $awesome]))
        outputs:
        - name: OUTPUT
          value: ($stdout)
        content: echo $GREETINGS
    - script:
        env:
        - name: INPUT
          value: ($OUTPUT)
        content: echo $INPUT
```

### Templating

Both `name` and `value` of an output can use [templating](./templating.md).
