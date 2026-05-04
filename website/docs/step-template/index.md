# Step template

A Chainsaw step template is very similar to a [test step](../step/index.md), except that it lives in its own resource and can be reused across tests.

Step templates help keep test suites DRY by factoring out common step patterns into a shared definition that multiple tests can invoke, each with their own arguments.

## Syntax

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: example
spec:
  # `bindings` defines default argument values for the template
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

!!! note
    The `try` field is required and must contain at least one operation.

### Reference

The full structure of `StepTemplate` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-StepTemplate).

### Example

The step template below applies and asserts a file whose name is determined by a binding.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: quick-start
spec:
  try:
  - apply:
      file: ($file)
  - assert:
      file: ($file)
```

## Lifecycle

The [same lifecycle](../step/index.md#lifecycle) as any test step applies to step templates.

## Bindings

Bindings defined in the template's `spec.bindings` act as **default values** for the template's parameters. When a test invokes the template via `with.bindings`, those bindings **override** the template's defaults.

This lets a template declare sensible defaults while still allowing callers to customise behaviour on a per-invocation basis.

### Example

The template below creates a ConfigMap whose name defaults to `from-template`. A calling test overrides that default by passing a different value.

**Template (`template.yaml`):**

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: template
spec:
  bindings:
  - name: input
    value: from-template   # default value
  try:
  - create:
      resource:
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ($input)
          namespace: ($namespace)
```

**Test:**

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - use:
      template: template.yaml
      with:
        bindings:
        - name: input
          value: from-test   # overrides the template default
```

## Invocation

To reference a step template from a test, use the `use` field in a step.

!!! note
    When `use` is specified, none of `try`, `catch`, `finally` and `cleanup` fields are allowed in the same step.

The full structure of `Use` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Use).

### Template

The path to the file containing the step template, **relative to the test file**.

### With

The `with` stanza contains the step template invocation arguments.

#### Bindings

Bindings registered here are merged with the template's own bindings before any operations execute. Test-provided bindings take precedence over same-named bindings defined in the template.

### Example

**Template (`step-template.yaml`):**

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: StepTemplate
metadata:
  name: quick-start
spec:
  try:
  - apply:
      file: ($file)
  - assert:
      file: ($file)
```

**Test:**

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - use:
      # path to the file containing the step template, relative to this test file
      template: step-template.yaml
      with:
        # bindings passed to the template; override any same-named template defaults
        bindings:
        - name: file
          value: configmap.yaml
```
