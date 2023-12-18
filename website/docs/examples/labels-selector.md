# Labels selector

Test steps:

1.  Creates a ConfigMap
1.  Asserts the ConfigMap contains the expected data

This test is similar to the [basic example](./basic.md) but Chainsaw is invoked with the `--selector` flag.

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
  name: labels-selector
  labels:
    # test labels
    test-suite: examples
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
# invoke chainsaw with `--selector` to filter tests to run
$ chainsaw test --selector test-suite=examples
```
