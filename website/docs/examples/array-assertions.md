# Array assertions

This example demonstrates how to perform complex assertions on arrays.

## Setup

See [Setup docs](./index.md#setup)

## Manifests

### `resources.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
  - name: container-1
    image: nginx-1
    env:
    - name: ENV_1
      value: value-1
  - name: container-2
    image: nginx-2
    env:
    - name: ENV_2
      value: value-2
  - name: container-3
    image: nginx-3
    env:
    - name: ENV_3
      value: value-3
```

### `assertions.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  # iterate over all containers having `name: container-1`
  ~.(containers[?name == 'container-1']):
    image: nginx-1
  # iterate over all containers, bind `$index` to the element index
  ~index.(containers):
    image: (join('-', ['nginx', to_string($index + `1`)]))
  # nested iteration
  ~index2.(containers):
    ~.(env):
      name: (join('_', ['ENV', to_string($index2 + `1`)]))
      value: (join('-', ['value', to_string($index2 + `1`)]))
```

## Test

### `chainsaw-test.yaml`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: array-assertions
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
chainsaw test
```
