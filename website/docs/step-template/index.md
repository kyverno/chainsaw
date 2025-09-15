# Step template

A Chainsaw step template is very similar to a test step, except that it lives in its own resource and can be reused across tests.

A step template can take leverage bindings to receive arguments from the calling test.

## Syntax

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: example
spec:
  # `bindings` defines arguments received from the calling test
  bindings: [...]
  # `try` defines operations to execute in the step
  try: [...]
  # `catch` defines operations to execute when the step fails
  catch: [...]
  # `finally` defines operations to execute at the end of the step
  finally: [...]
  # `cleanup` defines operations to execute at the end of the test
  cleanup: [...]
```

### Reference

The full structure of `StepTemplate` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-StepTemplate).

### Example

The step template below creates a config map with a `name` determined from a binding.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: example
spec:
  bindings:
  - name: input
    value: from-template
  try:
  - create:
      resource:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ($input)
```

## Lifecycle

The [same lifecycle](../step/index.md#lifecycle) as any test step applies to step templates.

## Invocation

To reference a step template from a test you need to specify the path to the file containing the template and eventually some bindings to simulate passing arguments to the step template.

!!! note
    When a step template is used, none of `try`, `catch`, `finally` and `cleanup` fields are allowed.

The full structure of `Use` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Use).

### Template

The path to the file containing the step template (relative to the test file).

### With

The `with` stanza contains the step template invocation details.

#### Bindings

Bindings can be registered prior to executing operations defined in the step template.
This effectively allows passing arguments to the step template.

### Example

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - use:
      # path to the file containing the step template
      template: template.yaml
      with:
        # bindings registered prior to executing operations defined in the step template
        # these bindings can be used to pass arguments to the step template
        bindings:
        - name: input
          value: from-test
```
