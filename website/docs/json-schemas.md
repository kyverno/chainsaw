# JSON schemas

JSON schemas for Chainsaw resources are available to enable validation and autocompletion in your IDE:

- [Configuration (v1alpha1)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/configuration-chainsaw-v1alpha1.json)
- [Configuration (v1alpha2)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/configuration-chainsaw-v1alpha2.json)
- [Test (v1alpha1)](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/test-chainsaw-v1alpha1.json)

## VS code

In VS code, simply add a comment on top of your YAML resources.

### Test

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: basic
spec:
  steps:
  - try:
    - apply:
        file: configmap.yaml
    - assert:
        file: configmap-assert.yaml
```

### Configuration

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/configuration-chainsaw-v1alpha2.json
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: congiguration
spec:
  fullName: true
  failFast: true
  forceTerminationGracePeriod: 5s
```

## Exporting schemas

Chainsaw can also export JSON schemas locally if you don't want to reference them from GitHub:

```bash
chainsaw export schemas <local path>
```

See [Chainsaw export schemas reference](./commands/chainsaw_export_schemas.md) for more details.
