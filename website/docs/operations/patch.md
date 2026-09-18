# Patch

The `patch` operation defines resources that should be modified in a Kubernetes cluster.

If the resource to be modified does not exist in the cluster, the step will fail.

## Patch semantics

The `patch` operation uses [JSON Merge Patch](https://www.rfc-editor.org/rfc/rfc7386) semantics. Object members are merged recursively, omitted object members generally remain unchanged, and setting an object member to `null` removes it. Arrays are replaced as complete values; they are not merged by Kubernetes list keys such as `name`.

For example, a patch containing only this Pod container:

```yaml
spec:
  containers:
    - name: app
      image: nginx:1.27
```

replaces the existing `spec.containers` array with that single object. Fields omitted from the replacement object, including defaulted or injected fields such as `resources` and `volumeMounts`, are not preserved in that array value. Kubernetes can therefore reject the update when those Pod fields are not allowed to change.

The native `patch` operation does not support RFC 6902 JSON Patch. To update a precise field inside an array, use the [`script` operation](./script.md) with a command such as `kubectl patch --type=json` and an explicit path like `/spec/containers/0/image`.

## Configuration

The full structure of `Patch` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Patch).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark: |
| [Outputs](../general/outputs.md) support           | :white_check_mark: |
| [Templating](../general/templating.md) support     | :white_check_mark: |
| [Operation checks](../general/checks.md) support   | :white_check_mark: |

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # use a specific file
        file: my-configmap.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # use glob pattern
        file: "configs/*.yaml"
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # use an URL
        file: https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/resource/valid.yaml
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        # specify resource inline
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: chainsaw-quick-start
          data:
            foo: bar
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example-subresource
spec:
  steps:
    - try:
        - patch:
            # specify resource inline, where a subresource is patched
            subresoruce: status
            resource:
              apiVersion: v1
              kind: Pod
              metadata:
                name: test-pod
              status:
                conditions:
                  - type: TestCondition
                    status: "True"
                    reason: TestReason
                    message: TestMessage
                    lastTransitionTime: "2026-03-28T10:51:35Z"
```

### Operation check

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - patch:
        file: my-configmap.yaml
        expect:
        - match:
            # this check applies only if the match
            # statement below evaluates to `true`
            apiVersion: v1
            kind: ConfigMap
          check:
            # an error is expected, this will:
            # - succeed if the operation failed
            # - fail if the operation succeeded
            ($error != null): true
```
