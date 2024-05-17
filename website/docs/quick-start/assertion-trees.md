# Use assertions trees

Whether verifying the scaling behavior of a database cluster or ensuring data consistency across instances, Chainsaw's assertion model provides the precision and flexibility needed for comprehensive validation.

Chainsaw allows declaring complex assertions with a simple and no-code approach, allowing assertions based on comparisons beyond simple equality, working with arrays, and other scenarios that could not be achieved before.

!!! tip
    Under the hood, Chainsaw uses [kyverno-json assertion trees](https://kyverno.github.io/kyverno-json/latest/intro/). Refer to the assertion trees documentation for more details on the supported syntax.

## Comparisons beyond simple equality

The assertion below will check that the number of replicas for a deployment is **greater than 3**.

Chainsaw doesn't need to know the exact expected number of replicas.
The `(replicas > 3)` expression will be evaluated until the assertion passes or the operation timeout expires (making the assertion fail).

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    # ...
    - assert:
        resource:
          apiVersion: v1
          kind: Deployment
          metadata:
            name: foo
          spec:
            (replicas > 3): true
    # ...
```

## Working with arrays

Chainsaw query language makes it easy to assert on arrays.
You can filter and transform arrays to select what you want to assert.

### Filtering 

In the example below we are creating a resource, then we assert that a condition with `type == 'Ready'` exists and has a field matching `status: 'True'`:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          spec:
            storage:
              secret:
                name: minio
                type: s3
            # ...
    - assert:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          status:
            # filter conditions array to keep elements where `type == 'Ready'`
            # and assert there's a single element matching the filter
            # and that this element status is `True`
            (conditions[?type == 'Ready']):
            - status: 'True'
```

## Comprehensive reporting

Chainsaw offers detailed resource diffs upon assertion failures. In the example below, the assertion failure message `metadata.annotations.foo: Invalid value: "null": Expected value: "bar"` is augmented with a resource diff.

It provides a clear view of discrepancies between expected and actual resources and gives more context around the specific failure (we can easily identify the owner of the offending pod for example).

```
| 09:55:50 | deployment | step-1   | ASSERT    | RUN   | v1/Pod @ chainsaw-rare-liger/*
| 09:56:20 | deployment | step-1   | ASSERT    | ERROR | v1/Pod @ chainsaw-rare-liger/*
    === ERROR
    ---------------------------------------------------
    v1/Pod/chainsaw-rare-liger/example-5477b4ff8c-tnhd9
    ---------------------------------------------------
    * metadata.annotations.foo: Invalid value: "null": Expected value: "bar"
    
    --- expected
    +++ actual
    @@ -1,10 +1,16 @@
      apiVersion: v1
      kind: Pod
      metadata:
    -  annotations:
    -    foo: bar
        labels:
          app: nginx
    +  name: example-5477b4ff8c-tnhd9
        namespace: chainsaw-rare-liger
    +  ownerReferences:
    +  - apiVersion: apps/v1
    +    blockOwnerDeletion: true
    +    controller: true
    +    kind: ReplicaSet
    +    name: example-5477b4ff8c
    +    uid: 118abe16-ec42-4894-83db-64479c4aac6f
      spec: {}
| 09:56:20 | deployment | step-1   | TRY       | DONE  |
```

## Next step

To continue our exploration of the main Chainsaw features, let's look at [bindings and resource templating](./resource-templating.md).
