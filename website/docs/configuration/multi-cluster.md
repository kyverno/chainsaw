# Multi cluster

Chainsaw supports testing against multiple clusters.

To use a specific cluster in a test (or test step) you will need to register the cluster either using the config file or using command line flags.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: custom-config
spec:
  # ...
  clusters:
    cluster-1:
      kubeconfig: /path/to/kubeconfig-1
      context: context-name
    cluster-2:
      kubeconfig: /path/to/kubeconfig-2
      context: context-name
  # ...
```

## Flag

The `--cluster` flag can appear multiple times and is expected to come in the following format `--cluster cluster-name=/path/to/kubeconfig[:context-name]`.

```bash
chainsaw test                                               \
    --cluster cluster-1=/path/to/kubeconfig-1               \
    --cluster cluster-2=/path/to/kubeconfig-2:context-name
```

## Usage

Once registered, a cluster can be specified at the `test`, `step` or `operation` level.
If specified at multiple levels, the most granular one is selected, effectively overriding the cluster specified at higher levels.

!!! example "Specifying cluster at the test level"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      # all steps and operations will be executed against the
      # cluster specified below (unless overridden)
      cluster: cluster-1
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

!!! example "Specifying cluster at the step level"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - # all operations will be executed against the
        # cluster specified below (unless overridden)
        cluster: cluster-1
        try:
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

!!! example "Specifying cluster at the operation level"

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
            # operation will be executed against the
            # cluster specified below
            cluster: cluster-1
            resource:
              apiVersion: v1
              kind: Deployment
              metadata:
                name: foo
              spec:
                (replicas > 3): true
        # ...
    ```

### `$client` binding

When a cluster is specified (whatever the `test`, `step` or `operation` level), the `$client` binding visible in operation expressions is always the Kubernetes client corresponding to the configured cluster.

In the example below, `$client` binding is a client to the configured `cluster-1`.

!!! example "`$client` binding"

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
            # operation will be executed against the
            # cluster specified below
            cluster: cluster-1
            resource:
              resource:
                (x_k8s_list($client, 'v1', 'Node')):
                    (length(@)): 1
        # ...
    ```
