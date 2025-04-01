# Use assertions

Chainsaw allows declaring complex assertions with a simple and no-code approach, allowing assertions based on comparisons beyond simple equality, working with arrays, and other scenarios that could not be achieved before.

!!! tip
    Under the hood, Chainsaw uses [kyverno-json assertion trees](https://kyverno.github.io/kyverno-json/latest/intro/). Refer to the assertion trees documentation for more details on the supported syntax.

## Basic assertion

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          apiVersion: apps/v1
          kind: Deployment
          metadata:
            name: coredns
            namespace: kube-system
          spec:
            replicas: 2
```

When asking Chainsaw to execute the assertion above, it will look for a deployment named `coredns` in the `kube-system` namespace and will compare the existing resource with the (partial) resource definition contained in the assertion.

In this specific case, if the field `spec.replicas` is set to 2 in the existing resource, the assertion will be considered valid.
If it is not equal to 2 the assertion will be considered failed.

This is the most basic assertion Chainsaw can evaluate.

## Slightly less basic assertion

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          apiVersion: apps/v1
          kind: Deployment
          metadata:
            labels:
              k8s-app: kube-dns
            namespace: kube-system
          spec:
            replicas: 2
```

This time we are not providing a resource name.

Chainsaw will look up **all** deployments with the `k8s-app: kube-dns` label in the `kube-system` namespace.
The assertion will be considered valid if **at least one** deployment matches the (partial) resource definition contained in the assertion.
If none match, the assertion will be considered failed.

Apart from the resource lookup process being a little bit more interesting, this kind of assertion is essentially the same as the previous one.
Chainsaw is basically making a decision by comparing an actual and expected resource.

## Beyond simple equality

The assertion below will check that the number of replicas for a deployment is **greater than 1 AND less than 4**.

Chainsaw doesn't need to know the exact expected number of replicas.
The `(replicas > 1 && replicas < 4)` expression will be evaluated until the result is `true` or the operation timeout expires (making the assertion fail).

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          apiVersion: apps/v1
          kind: Deployment
          metadata:
            name: coredns
            namespace: kube-system
          spec:
            (replicas > `1` && replicas < `4`): true
```

!!! tip
    To indicate that a key or value in the YAML document is an expression, simply place the element between parenthesis:

    - `this is an expression` -> interpreted as a `string`
    - `(this is an expression)` -> interpreted as a JMESPath expression


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
            ...
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

### Iterating

Being able to filter arrays allows selecting the elements to be processed.

On top of that, Chainsaw allows iterating over array elements to validate each item separately.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          apiVersion: apps/v1
          kind: Deployment
          metadata:
            labels:
              k8s-app: kube-dns
            namespace: kube-system
          spec:
            template:
              spec:
                # the `~` modifier tells Chainsaw to iterate over the array elements
                ~.(containers):
                  securityContext: {}
```

This assertion uses the `~` modifier and Chainsaw will evaluate descendants once per element in the array.

## Comprehensive reporting

Chainsaw offers detailed resource diffs upon assertion failures.

In the example below, the assertion failure message (`metadata.annotations.foo: Invalid value: "null": Expected value: "bar"`) is augmented with a resource diff.

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

To continue our exploration of the main Chainsaw features, let's look at [bindings](./bindings.md) and [resource templating](./resource-templating.md) next.
