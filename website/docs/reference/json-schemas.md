# JSON schemas

JSON schemas for Chainsaw are available:

- [Configuration (v1alpha1)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/configuration-chainsaw-v1alpha1.json)
- [Configuration (v1alpha2)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/configuration-chainsaw-v1alpha2.json)
- [Test (v1alpha1)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/test-chainsaw-v1alpha1.json)
- [StepTemplate (v1alpha1)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/test-chainsaw-v1alpha1.json)

They can be used to enable validation and autocompletion in your IDE.

## VS code

In VS code, simply add a comment on top of your YAML resources.

### Test

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        file: configmap.yaml
    - assert:
        file: configmap-assert.yaml
```

### StepTemplate

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/steptemplate-chainsaw-v1alpha1.json
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

### Configuration

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/configuration-chainsaw-v1alpha2.json
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  timeouts:
    apply: 45s
    assert: 20s
    cleanup: 45s
    delete: 25s
    error: 10s
    exec: 45s
  cleanup:
    skipDelete: false
  execution:
    failFast: true
    parallel: 4
```

## Exporting schemas

Chainsaw can also export JSON schemas locally if you don't want to reference them from GitHub:

```bash
chainsaw export schemas <local path>
```

See [chainsaw export schemas](./commands/chainsaw_export_schemas.md) command documentation for more details.
