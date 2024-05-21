# Control your cleanup

Unless configured differently, by default Chainsaw will **automatically remove the resources it created** after a test finishes.

Cleanup happens in reverse order of creation (created last, cleaned up first).
This is important, especially when the controller being tested makes use of `finalizers`.

!!! tip "Overriding cleanup timeout"
    Note that Chainsaw performs a blocking deletion, that is, it will wait until the resource is not present anymore in the cluster before proceeding with the next resource cleanup.

## Timeout

A global cleanup timeout can be defined at the configuration level or using command line flags.

It can also be overridden on a per-test or per-step basis but not at the operation level.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  timeouts:
    # cleanup timeout at the test level
    cleanup: 30s
  steps:
  - timeouts:
      # cleanup timeout at the step level
      cleanup: 2m
    try: ...
```

## Automatic cleanup

After a test, every resource created by Chainsaw will be automatically deleted. This applies to `create` and `apply` operations.

In the logs below we can see Chainsaw deletes the previously created resource:

```
    | 15:21:29 | quick-start | @setup   | CREATE    | OK    | v1/Namespace @ chainsaw-cute-cod
    | 15:21:29 | quick-start | step-1   | TRY       | RUN   |
    | 15:21:29 | quick-start | step-1   | APPLY     | RUN   | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | CREATE    | OK    | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | APPLY     | DONE  | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | ASSERT    | RUN   | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | ASSERT    | DONE  | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | TRY       | DONE  |
    === step cleanup process start ===
    | 15:21:29 | quick-start | step-1   | CLEANUP   | RUN   |
    | 15:21:29 | quick-start | step-1   | DELETE    | RUN   | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | DELETE    | OK    | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | DELETE    | DONE  | v1/ConfigMap @ chainsaw-cute-cod/chainsaw-quick-start
    | 15:21:29 | quick-start | step-1   | CLEANUP   | DONE  |
    === step cleanup process end ===
    === test cleanup process start ===
    | 15:21:29 | quick-start | @cleanup | DELETE    | RUN   | v1/Namespace @ chainsaw-cute-cod
    | 15:21:29 | quick-start | @cleanup | DELETE    | OK    | v1/Namespace @ chainsaw-cute-cod
    | 15:21:34 | quick-start | @cleanup | DELETE    | DONE  | v1/Namespace @ chainsaw-cute-cod
    === test cleanup process end ===
```

## Manual cleanup

Under certain circumstances, automatic cleanup is not enough and we want to execute custom operations.

Chainsaw allows registering cleanup operations that will be run after automatic cleanup.
Custom cleanup operations live at the test step level:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
    # this step will create a local cluster
  - try:
    - script:
        timeout: 1m
        content: |
          kind create cluster --name dynamic --kubeconfig ./dynamic
    # at cleanup time, we want to delete the local cluster we created
    # and remove the associated kubeconfig
    cleanup:
    - script:
        content: |
          kind delete cluster --name dynamic
    - script:
        content: |
          rm -f ./dynamic
```

## Next step

At this point, we covered the main Chainsaw features.

Look at the [next steps](./next-steps.md) section to find out what to do next.
