# Bindings

You can think of bindings as a side context where you can store and retrieve data by name.

This is particularly useful when some data is only known at runtime.
For example, to pass data from one operation to another, to implement resource templating, to fetch data from an external system, or anything that needs to be computed at runtime.

## Syntax

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

### Reference

Browse the [reference documentation](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Binding) to see the syntax details and where bindings can be declared.

### Inheritance

Bindings can be configured at the test, step or operation level.

All bindings configured at a given level are automatically [inherited](./inheritance.md) at lower levels.

### Immutability

Bindings are immutable. This means two bindings can have the same name without overwriting each other.

When a binding is registered it potentially hides other bindings with the same name.

When this binding goes out of scope, previously registered bindings with the same name become visible again.

### Templating

Both `name` and `value` of a binding can use [templating](./templating.md).

## Built-in bindings

Chainsaw offers some built-in bindings you can directly use in your tests, steps and operations.

Browse the [built-in bindings list](../reference/builtins.md) to find available bindings.

## Binding objects with `x_k8s_get()` to refer/assert against

```yaml
    - description: >-
        1. Look up some details about a Deployment object with 'x_k8s_get(...)', and store them in a binding.
        2. Make use of the binding in a `script` block through environment variables.
        3. Refer to the same bound Deployment object in an `assert` block.
      bindings:
        - name: xpNS
          value: crossplane-system
        - name: xpDeploy
          # Arguments to 'x_k8s_get(any, string, string, string, string)':
          #
          # 1. any: Supply the '$client' built-in binding here, which is needed to connect to the cluster:
          #         https://kyverno.github.io/chainsaw/latest/reference/builtins/
          #
          # 2. string: 'apiVersion' field on the object.
          #
          # 3. string: 'Kind' field on the object.
          #
          # 4. string: The namespace of the object. If the object type is not namespaced, this field can be an empty
          #            string, or any string; it doesn't seem to matter.
          #
          # 5. string: The name of the object.
          value: (x_k8s_get($client, 'apps/v1', 'Deployment', $xpNS, 'crossplane'))
      try:
        - script:
            bindings:
              - # Re-bind the version label from the Deployment to give us a more succinct name to refer to.
                name: deployVersion
                # If the label key has any periods in it, double-quote the whole key.
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
            env:
              - # Refer to the Deployment version label the long way.
                name: DEP_VER_LONG
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
              - # Refer to the Deployment version label the short way, through the additional binding scoped to this
                # 'script' block.
                name: DEP_VER_SHORT
                value: ($deployVersion)
            # The version values printed by the script here will be the same, even though they were derived through
            # slightly different routes.
            content: |-
              echo "DEP_VER_LONG:  '$DEP_VER_LONG'"
              echo "DEP_VER_SHORT: '$DEP_VER_SHORT'"
        - assert:
            bindings:
              - name: depVer
                value: ($xpDeploy.metadata.labels."app.kubernetes.io/version")
            resource:
              apiVersion: v1
              kind: Pod
              metadata:
                # Match the pod(s) with namespace and label selectors.
                namespace: ($xpNS)
                labels:
                  app: crossplane
              # Make a binding that holds this Pod's version from its label.
              (metadata.labels."app.kubernetes.io/version")->podVer:
                # Assert that the version labels for the Deployment and this Pod equal each other, using the
                # 'semver_compare()' function described here:
                # https://kyverno.io/docs/writing-policies/jmespath/#semver_compare
                (semver_compare($depVer, $podVer)): true
```
