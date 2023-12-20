# Assertion trees

Test steps:

1.  Creates a Deployment with 2 replicas
1.  Asserts that the number of replicas is > 1 and that `status.replicas == spec.replicas`

## Setup

See [Setup docs](./index.md#setup)

## Manifests

### `resources.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
```

### `assertions.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
(status.replicas == spec.replicas): true
spec:
  (replicas > `1`): true
```

## Test

### `chainsaw-test.yaml`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: assertion-tree
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
