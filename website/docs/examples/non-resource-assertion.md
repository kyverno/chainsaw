# Non resource assertion

This test example demonstrates how to perform assertions not based on resources.

1.  Asserts that the number of nodes in the cluster is equal to 1

## Setup

See [Setup docs](./index.md#setup)

## Manifests

### `assertions.yaml`

```yaml
(x_k8s_list($client, 'v1', 'Node')):
  (length(items): 1
```

## Test

### `chainsaw-test.yaml`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: non-resource-assertion
spec:
  steps:
  - try:
    - assert:
        file: assertions.yaml
```

## Execute

```bash
chainsaw test
```
