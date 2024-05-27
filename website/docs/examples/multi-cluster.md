# Multi-cluster setup

Chainsaw supports registering and using [multiple clusters](../configuration/options/clusters.md) in tests.

We can also register clusters dynamically and combine this with cluster selection to achieve scenarios where clusters are dynamically allocated in a test step, used in the following steps, and cleaned up at the end.

The following test demonstrates such a scenario by creating a local kind cluster in the first, using it in the second step, and configuring a [cleanup](../quick-start/cleanup.md) script to delete the cluster when the test terminates:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    # create a local cluster
    - script:
        timeout: 1m
        content: |
          kind create cluster --name dynamic --kubeconfig ./dynamic
    # register `cleanup` operations to delete the cluster
    # at the end of the test
    cleanup:
    - script:
        content: |
          kind delete cluster --name dynamic
    - script:
        content: |
          rm -f ./dynamic
    # register the `dynamic` cluster in this step
  - clusters:
      dynamic:
        kubeconfig: ./dynamic
    # and use the `dynamic` cluster for all operations in the step
    cluster: dynamic
    try:
    - apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
            namespace: default
          data:
            foo: bar
    - assert:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
            namespace: default
          data:
            foo: bar
```

Running the test above will produce the following output:

```
    | 10:44:53 | example | @setup   | CREATE    | OK    | v1/Namespace @ chainsaw-useful-seahorse
    | 10:44:53 | example | step-1   | TRY       | RUN   |
    | 10:44:53 | example | step-1   | SCRIPT    | RUN   |
        === COMMAND
        /bin/sh -c kind create cluster --name dynamic --kubeconfig ./dynamic
    | 10:45:10 | example | step-1   | SCRIPT    | LOG   |
        === STDERR
        Creating cluster "dynamic" ...
         • Ensuring node image (kindest/node:v1.27.3) 🖼  ...
         ✓ Ensuring node image (kindest/node:v1.27.3) 🖼
         • Preparing nodes 📦   ...
         ✓ Preparing nodes 📦 
         • Writing configuration 📜  ...
         ✓ Writing configuration 📜
         • Starting control-plane 🕹️  ...
         ✓ Starting control-plane 🕹️
         • Installing CNI 🔌  ...
         ✓ Installing CNI 🔌
         • Installing StorageClass 💾  ...
         ✓ Installing StorageClass 💾
        Set kubectl context to "kind-dynamic"
        You can now use your cluster with:
        
        kubectl cluster-info --context kind-dynamic --kubeconfig ./dynamic
        
        Thanks for using kind! 😊
    | 10:45:10 | example | step-1   | SCRIPT    | DONE  |
    | 10:45:10 | example | step-1   | TRY       | DONE  |
    | 10:45:10 | example | step-2   | TRY       | RUN   |
    | 10:45:10 | example | step-2   | APPLY     | RUN   | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | CREATE    | OK    | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | APPLY     | DONE  | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | ASSERT    | RUN   | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | ASSERT    | DONE  | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | TRY       | DONE  |
    | 10:45:10 | example | step-2   | CLEANUP   | RUN   |
    | 10:45:10 | example | step-2   | DELETE    | RUN   | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | DELETE    | OK    | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | DELETE    | DONE  | v1/ConfigMap @ default/quick-start
    | 10:45:10 | example | step-2   | CLEANUP   | DONE  |
    | 10:45:10 | example | step-1   | CLEANUP   | RUN   |
    | 10:45:10 | example | step-1   | SCRIPT    | RUN   |
        === COMMAND
        /bin/sh -c kind delete cluster --name dynamic
    | 10:45:10 | example | step-1   | SCRIPT    | LOG   |
        === STDERR
        Deleting cluster "dynamic" ...
        Deleted nodes: ["dynamic-control-plane"]
    | 10:45:10 | example | step-1   | SCRIPT    | DONE  |
    | 10:45:10 | example | step-1   | SCRIPT    | RUN   |
        === COMMAND
        /bin/sh -c rm -f ./dynamic
    | 10:45:10 | example | step-1   | SCRIPT    | DONE  |
    | 10:45:10 | example | step-1   | CLEANUP   | DONE  |
    | 10:45:10 | example | @cleanup | DELETE    | RUN   | v1/Namespace @ chainsaw-useful-seahorse
    | 10:45:11 | example | @cleanup | DELETE    | OK    | v1/Namespace @ chainsaw-useful-seahorse
    | 10:45:16 | example | @cleanup | DELETE    | DONE  | v1/Namespace @ chainsaw-useful-seahorse
```