# Negative testing

Negative testing is the process of testing cases that are supposed to fail. That is, a test expects errors to happen and if the expected errors don't occur the test must fail.

Chainsaw supports negative testing by letting you decide what should be considered an error or not.

!!! tip
    By default, Chainsaw will consider an operation failed if there was an error executing it (non-zero exit code in scripts and commands, error returned by the API server when calling into Kubernetes, etc...).

## Script case

The test below expects an error and validates the returned error message:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - script:
        content: kubectl get foo
        check:
          ($error != null): true
          ($stderr): |-
            error: the server doesn't have a resource type "foo"
```

If for whatever reason, the `kubectl get foo` doesn't return an error, or the message received in standard error output is not `error: the server doesn't have a resource type "foo"`, Chainsaw will consider the operation failed.

If it returns an error and the expected error message, Chainsaw will consider the operation successful.

## Working with resources

The test below tries to apply resources in a cluster but expects the operation to fail:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        file: resources.yaml
        expect:
          # check that applying the resource failed
        - check:
            ($error != null): true
```

If applying the resource succeeded, Chainsaw will consider the operation failed.

On the other hand, if applying the resource fails, Chainsaw will consider the operation to be successful.

### Resource matching

In the previous example, if the `resources.yaml` contains multiple resources, but only some of them may be expected to fail.

Chainsaw allows matching resources when evaluating checks:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        file: resources.yaml
        expect:
          # the check below only applies if the resource being checked
          # matches the condition defined in the `match` field
        - match:
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: quick-start
          check:
            ($error != null): true
```

Using the `match` field, we can easily target failures related to specific resources.