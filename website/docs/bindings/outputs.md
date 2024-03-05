# Outputs

Binding outputs can be useful to communicate and reuse computation results between operations.

## Supported operations

Currently, only `script` and `command` operations support outputs.

## Lifetime of outputs

Once an output has been added in the form of a binding, this binding will be available to all following operations **in the same step**.

Currently, outputs do not cross the step boundaries.

## Examples

The example below defines invokes a `kubectl` command to get a configmap from the cluster in json format.

The json output is then parsed and added to the `$cm` binding and the next operation performs an assertion on it by reading the binding instead of querying the cluster.

!!! example "Output in script"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      # ...
      - try:
        - script:
            content: kubectl get cm quick-start -n $NAMESPACE -o json
            outputs:
            - name: cm
              value: (json_parse($stdout))
        - assert:
            resource:
              ($cm):
                metadata:
                  (uid != null): true
      # ...
    ```

!!! example "Output in command"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      # ...
      - try:
        - command:
            entrypoint: kubectl
            args:
            - get
            - cm
            - quick-start
            - -n
            - $NAMESPACE
            - -o
            - json
            outputs:
            - name: cm
              value: (json_parse($stdout))
        - assert:
            resource:
              ($cm):
                metadata:
                  (uid != null): true
      # ...
    ```
