# Command

The `command` operation provides a mean to execute a specific command during the test step.

!!! warning
    Command arguments are not going through shell expansion.
    
    It's crucial to consider potential differences in behavior, as Chainsaw may interpret them differently compared to regular shell environments.

## Configuration

The full structure of the `Command` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Command).

### Features

| Supported features                                 |                    |
|----------------------------------------------------|:------------------:|
| [Bindings](../general/bindings.md) support         | :white_check_mark: |
| [Outputs](../general/outputs.md) support           | :white_check_mark: |
| [Templating](../general/templating.md) support     | :x:                |
| [Operation checks](../general/checks.md) support   | :white_check_mark: |

### KUBECONFIG

- Unless `--no-cluster` is specified, Chainsaw always executes commands in the context of a temporary `KUBECONFIG`, built from the configured target cluster.
- This specific `KUBECONFIG` has a single cluster, auth info and context configured (all named `chainsaw`).

### Environment variables expansion

Chainsaw will expand environment variables in the form of `$VARIABLE_NAME`. If you need to provide the `$` sign you can do it by escaping it with `$$`.

This matches [Kubernetes' behavior](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#entrypoint) for container command and args fields.

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - command:
        entrypoint: echo
        args:
        - hello chainsaw
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
    - command:
        entrypoint: echo
        args:
        - hello chainsaw
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
```
