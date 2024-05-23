# Bindings

Chainsaw has a concept of `bindings` which can be seen as an execution context.

Bindings are used in [assertion trees](../../operations/check.md) and [resource templating](../../configuration/options/templating.md), as well as when using the `--values` flag when invoking Chainsaw.

While some bindings are built-in and provided by Chainsaw, it's possible to define custom bindings at the test, step, or individual operation level.
Those bindings can in turn be used to create custom environment variables in `command` and `script` operations.

## Immutability

It's important to note that bindings are **immutable**, it's not possible to overwrite a binding and **two bindings with the same name can exist without overwriting each other**.

When a binding value is resolved, Chainsaw will walk the binding chain upwards until it finds a binding with the expected name.
Effectively, the last one registered in the chain will be used.

## Templating

A binding `name` supports templating.

The name of a binding can therefore be dynamic and depend on [values](../../quick-start/advanced/values.md) or other bindings.

## Usage

The example below defines custom bindings at the test level.

!!! example "Test level bindings"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      # bindings defined at the test level are available to all steps and operations
      bindings:
      - name: hello
        value: hello
      - name: chainsaw
        value: chainsaw
      # `greetings` depends on `$hello` and `$chainsaw` bindings defined above
      - name: greetings
        value: (join(' ', [$hello, $chainsaw]))
      steps:
      - try:
        - script:
            # custom environment variables, defined using custom bindings
            env:
            - name: GREETINGS
              value: ($greetings)
            content: echo $GREETINGS
            check:
              ($error): ~
              ($stdout): hello chainsaw
    ```

The example below is similar to the previous one but also defines custom bindings at the step level.

!!! example "Step level bindings"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      # bindings defined at the test level are available to all steps and operations
      bindings:
      - name: hello
        value: hello
      - name: chainsaw
        value: chainsaw
      steps:
      - bindings:
        # `greetings` depends on `$hello` and `$chainsaw` bindings defined at higher levels
        - name: greetings
          value: (join(' ', [$hello, $chainsaw]))
        try:
        - script:
            # custom environment variables, defined using custom bindings
            env:
            - name: GREETINGS
              value: ($greetings)
            content: echo $GREETINGS
            check:
              ($error): ~
              ($stdout): hello chainsaw
    ```

Finally, custom bindings can also be defined at the operation level.

!!! example "Operation level bindings"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      # bindings defined at the test level are available to all steps and operations
      bindings:
      - name: hello
        value: hello
      steps:
      - bindings:
        - name: chainsaw
          value: chainsaw
        try:
        - script:
            bindings:
            # `greetings` depends on `$hello` and `$chainsaw` bindings defined at the higher levels
            - name: greetings
              value: (join(' ', [$hello, $chainsaw]))
            # custom environment variables, defined using custom bindings
            env:
            - name: GREETINGS
              value: ($greetings)
            content: echo $GREETINGS
            check:
              ($error): ~
              ($stdout): hello chainsaw
    ```

## Outputs

Under certain conditions, bindings can also be used to add computed results to the context.

See [Outputs](./outputs.md) for details.
