# Commands

In addition to collecting pod logs and events, Chainsaw also supports arbitrary commands as collectors.

## Configuration

The full structure of the `Command` resource is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Command).

### Simple command

A `Command` must have at least the `entrypoint` defined.

!!! example "Simple command"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - command:
            entrypoint: time
    ```

### Command with arguments

`Command` arguments can be provided using the `args` field.

!!! example "Command arguments"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - command:
            entrypoint: kubectl
            args:
            - get
            - pod
            - -n
            - $NAMESPACE
    ```

### Timeout

An optional `timeout` can be configured.

!!! note
    Note that the `timeout` lives at the operation level, not at the `script` level.

!!! example "Timeout example"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - command:
            entrypoint: kubectl
            args:
            - get
            - pod
            - -n
            - $NAMESPACE
          timeout: 30s
    ```
