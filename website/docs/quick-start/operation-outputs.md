# Use operation outputs

Operation outputs can be useful for communicating and reusing computation results across operations.

## Lifetime of outputs

Once an output has been added to the bindings context, this binding will be available to all following operations **in the same step**.

Currently, outputs do not cross the step boundaries.

## Matching

An output supports an optional `match` field. The `match` is used to conditionally create a binding.

In the case of applying a file, for example, the file may contain multiple resources. The `match` can be used to select the resource to use for creating the binding.

## Load an existing resource

The example below invokes a `kubectl` command to get a configmap from the cluster in json format.

The json output is then parsed and added to the `$cm` binding and the next operation performs an assertion on it by reading the binding instead of querying the cluster.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        content: kubectl get cm quick-start -n $NAMESPACE -o json
        outputs:
          # parse stdout json output and bind the result to `$cm`
        - name: cm
          value: (json_parse($stdout))
    - assert:
        resource:
          ($cm):
            metadata:
              (uid != null): true
```

## Match a resource

The example below applies resources from a file.

When the resource being applied is a configmap, we bind the resource to an output to print its UID in the next operation.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        file: ./resources.yaml
        outputs:
          # match the configmap resource and bind it to `$cm`
        - match:
            apiVersion: v1
            kind: ConfigMap
          name: cm
          value: (@)
    - script:
        env:
        - name: UID
          value: ($cm.metadata.uid)
        content: echo $UID
```

## Next step

In the next section, we will look at the three main elements of a test step, [the `try`, `catch` and `finally` blocks](./try-catch.md).
