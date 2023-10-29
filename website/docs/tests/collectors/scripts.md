# Scripts

In addition to collecting pod logs and events, Chainsaw also supports arbitrary scripts as collectors.

!!!warning "Shell"

    Scripts require a shell to run, Chainsaw executes scripts with `sh -c ...`.

    If no shell is available, `Scripts` can't be used.

## Configuration

The full structure of the `Script` resource is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

A `Script` must have a `content` defined.

!!! example "Simple command"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - exec:
              script:
                content: |
                  echo "an error has occured"
    ```

### Timeout

An optional `timeout` can be configured.

!!! note
    Note that the `timeout` lives at the `exec` level, not at the `script` level.

!!! example "Timeout example"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - exec:
              timeout: 30s
              script:
                content: |
                  echo "an error has occured"
    ```
