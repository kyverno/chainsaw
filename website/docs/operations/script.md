# Script

The `script` operation provides a means to run a script during the test step.


!!! tip
    By default, the shell used to run the script is `sh`. You can specify a custom shell and shell arguments to override what and how the shell is invoked.

## Configuration

The full structure of `Script` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

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

### Working directory

Use `workDir` to set the working directory in which the script will be executed. When omitted, the script runs in the directory containing the test file.

`workDir` accepts either an absolute path or a relative path. A relative path is resolved against the directory containing the test file.

```yaml
- script:
    workDir: /path/to/project
    content: |
      make build
```

### Log suppression

Two fields control whether script output appears in the test logs:

- `skipLogOutput`: suppresses the script's stdout/stderr output from the logs. Useful for reducing noise or hiding sensitive values.
- `skipCommandOutput`: suppresses the shell invocation itself from the logs.

```yaml
- script:
    skipLogOutput: true
    skipCommandOutput: true
    content: |
      echo secret-value
```

## Examples

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        content: |
          echo "hello chainsaw"
```

### Custom shell and shell arguments

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        # use `bash` shell
        shell: bash
        # invoke `bash` with `-c`
        shellArgs:
        - -c
        content: |
          echo "hello chainsaw"
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
    - script:
        content: |
          echo "hello chainsaw"
        check:
          # an error is expected, this will:
          # - succeed if the operation failed
          # - fail if the operation succeeded
          ($error != null): true
```
