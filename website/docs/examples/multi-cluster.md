# Multi-cluster setup

Chainsaw supports testing against multiple clusters.

To use a specific cluster in a test (or step) you will need to register the cluster either using the config file or using command line flags.

Since `v0.2.1` you can also register clusters dynamically at the test, step and operation levels. This is particularly useful when a cluster is created in a test step and used in subsequent steps.

## Register clusters

### In Configuration

Chainsaw configuration file has a `clusters` field you can use to register existing clusters. Clusters are registered by name, and point to a kubeconfig file and an optional context name.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  clusters:
    # this cluster will use the default (current) context
    # configured in the kubeconfig file
    cluster-1:
      kubeconfig: /path/to/kubeconfig-1
    # this cluster will use the context named `context-2`
    # in the kubeconfig file
    cluster-2:
      kubeconfig: /path/to/kubeconfig-2
      context: context-2
```

### Using flags

The `--cluster` flag can appear multiple times and is expected to come in the following format `--cluster cluster-name=/path/to/kubeconfig[:context-name]`.

```bash
chainsaw test                                               \
    --cluster cluster-1=/path/to/kubeconfig-1               \
    --cluster cluster-2=/path/to/kubeconfig-2:context-2
```

Using the flags above is equivalent to the cluster registration using a configuration file from the previous example.

!!! tip "Precedence"
    Remember that **flags take precedence over the configuration file when** both are specified.


### Dynamic registration

The test below illustrates dynamic cluster registration:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  # register clusters at the test level
  # those clusters will be inherited in all steps and operations
  # and can be overridden
  clusters:
    cluster-1:
      kubeconfig: /path/to/kubeconfig-1
    cluster-2:
      kubeconfig: /path/to/kubeconfig-2
      context: context-2
  steps:
  - clusters:
      # register clusters at the step level
      # `cluster-1` will be overridden for this particular step
      cluster-1:
        kubeconfig: /path/to/another-kubeconfig-1
    try:
      # operation runs against `cluster-1`
    - cluster: cluster-1
      apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
            namespace: default
          data:
            foo: bar
  - try:
    - clusters:
        # register clusters at the operation level
        # `cluster-1` will be overridden for this particular operation
        cluster-1:
          kubeconfig: /path/to/yet-another-kubeconfig-1
      # operation runs against `cluster-1`
      cluster: cluster-1
      apply:
        resource:
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: quick-start
            namespace: default
          data:
            foo: bar
```

!!! tip
    The default cluster uses the `''` name. You can override it with:

    ```yaml
    clusters:
      '':
        kubeconfig: /path/to/kubeconfig
        context: context
    ```

## Use clusters

Registers clusters can be used by name and assigned/overridden at the test, step or individual operation level using the `cluster` field.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  # assigning a cluster at the test level
  # `cluster-1` will be used in all steps and operations
  # (if not overridden)
  cluster: cluster-1
  steps: ...
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
    # assigning a cluster at the step level
    # `cluster-1` will be used in all step operations
    # (if not overridden)
  - cluster: cluster-1
    try: ...
    # `cluster-1` will be used in all step operations
    # (if not overridden)
  - cluster: cluster-2
    try: ...
---
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
      # assigning a cluster at the operation level
    - cluster: cluster-1
      apply: ...
      # assigning a cluster at the operation level
    - cluster: cluster-2
      apply: ...
```

## Combine both

It is completely possible to combine both dynamic registration and cluster selection to achieve scenarios where clusters are dynamically allocated in a test step, used in the following steps, and cleaned up at the end.

The following test demonstrates such a scenario by creating a local kind cluster in the first, using it in the second step, and configuring a `cleanup` script to delete the cluster when the test terminates:

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