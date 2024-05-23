# Outputs

Binding outputs can be useful to communicate and reuse computation results between operations.

## Supported operations

Currently, only `script` and `command` operations support outputs.

## Lifetime of outputs

Once an output has been added in the form of a binding, this binding will be available to all following operations **in the same step**.

Currently, outputs do not cross the step boundaries.

## Matching

An output supports an optional `match` field. The `match` is used to conditionally create a binding.

In the case of applying a file, for example, the file may contain multiple resources. The `match` can be used to select the resource to use for creating the binding.

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
            - match:
                (json_parse($stdout)):
                  apiVersion: v1
                  kind: ConfigMap
              name: cm
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
            - match:
                (json_parse($stdout)):
                  apiVersion: v1
                  kind: ConfigMap
              name: cm
              value: (json_parse($stdout))
        - assert:
            resource:
              ($cm):
                metadata:
                  (uid != null): true
      # ...
    ```
