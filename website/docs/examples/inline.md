# Inline resources

Test steps:

1.  Creates a ConfigMap
1.  Asserts the ConfigMap contains the expected data

## Setup

See [Setup docs](./index.md#setup)

## Test

### `chainsaw-test.yaml`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: inline-resources
spec:
  steps:
  - try:
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
    - assert:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
          data:
            foo: bar
```

## Execute

```bash
chainsaw test
```
