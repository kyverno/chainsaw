# JSON schemas

JSON schemas for Chainsaw resources are available to enable validation and autocompletion in your IDE:

- [Configuration](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/configuration-chainsaw-v1alpha1.json)
- [Test](https://github.com/kyverno/chainsaw/blob/main/.schemas/json/test-chainsaw-v1alpha1.json)

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

## Applying CRDs

An alternative is to apply the Chainsaw CRDs in a kubernetes cluster and let the Kubernetes extension do the rest.

CRD definitions are provided in our [GitHub repository](https://github.com/kyverno/chainsaw/tree/main/config/crds).

The command below will apply the Chainsaw CRDs to the configured cluster:

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/chainsaw/main/config/crds/chainsaw.kyverno.io_configurations.yaml
kubectl apply -f https://raw.githubusercontent.com/kyverno/chainsaw/main/config/crds/chainsaw.kyverno.io_tests.yaml
```

## Exporting schemas

Chainsaw can export JSON schemas locally.

```bash
chainsaw export schemas <local path>
```

See [Chainsaw export schemas reference](./commands/chainsaw_export_schemas.md) for more details.
