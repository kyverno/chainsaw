# Use try, catch and finally

A test step is made of 3 main blocks used to determine the actions Chainsaw will perform when executing the step, depending on the test outcome.

- The [try](../steps/try.md) block *(required)*
- The [catch](../steps/catch.md) block *(optional)*
- The [finally](../steps/finally.md) block *(optional)*

Operations defined in the `try` block are executed first, then:

- If an operation fails to execute, Chainsaw won't execute the remaining operations and will execute **all** operations defined in the `catch` block instead (if any).
- If all operations succeed, Chainsaw will NOT execute operations defined in the `catch` block (if any).
- Regardless of the step outcome (success or failure), Chainsaw will execute **all** operations defined in the `finally` block (if any).

!!! note

    Note that all operations coming from the `catch` or `finally` blocks are executed. If one operation fails, Chainsaw will mark the test as failed and continue executing with the next operation.

## Cleanup

At the end of a test, Chainsaw automatically cleans up the resources created during the test (cleanup is done in the opposite order of creation).

All operations from the `catch` and `finally` blocks are executed before the cleanup process kicks in. This order allows analyzing the resources that potentially caused the step failure before they are deleted.

## Catch

Operations in a `catch` block are executed only when the step is considered failed.

This is particularly useful to collect additional information to help understand what caused the failure.

In the example below, the test contains a `catch` block to collect events in the cluster when an operation fails in the step.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # ...
    - assert:
        # ...
    # collect events in the `catch` block
    # will be executed only if an operation failed
    catch:
    - events: {}
```

## Finally

Operations in a `finally` block will always execute regardless of the success or failure of the test step.

This is particularly useful to perform manual cleanup.

In the example below we create a local cluster in a script operation. The cluster deletion script is added to the `finally` block, guaranteeing the cluster will be deleted regardless of the test outcome.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
    # create a local cluster
  - try:
    - script:
        timeout: 1m
        content: |
          kind create cluster --name dynamic --kubeconfig ./dynamic
    - apply:
        # ...
    - assert:
        # ...
    # add cluster deletion script in the `finally` block
    # to guarantee the cluster will be deleted after the test
    finally:
    - script:
        content: |
          kind delete cluster --name dynamic
    - script:
        content: |
          rm -f ./dynamic
```

## Next step

Every operation in a test must be executed in a timely fashion. In the next section, we will see how you can [control your timeouts](./timeouts.md).
