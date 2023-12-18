# Basic

Test steps:

1.  Creates a ConfigMap
1.  Asserts the ConfigMap contains the expected data

## Setup

See [Setup docs](./index.md#setup)

## Manifests

### `resources.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: quick-start
data:
  foo: bar
```

### `assertions.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: quick-start
data:
  foo: bar
```

## Test

### `chainsaw-test.yaml`

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
        file: resources.yaml
    - assert:
        file: assertions.yaml
```

## Execute

```bash
$ chainsaw test
```
