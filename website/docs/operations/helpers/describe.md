# Describe

Describing resources present in the cluster can help understand what happened and troubleshoot test failures.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Describe` resource is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Describe).

!!! warning "Deprecated syntax"
    You can specify the `resource` directly instead of using `apiVersion` and `kind`.
    
    **This is a deprecated syntax though and will be removed in a future version.**

### Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

### All namespaces

When used with a namespaced resource, it is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

## Usage examples

### Describe pods

If a `name` is specified, Chainsaw will describe resources that have a name starting with the provided `name` in the test namespace (or in the cluster if it is a clustered-level resource).

!!! example "Describe pods in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        catch:
        - describe:
            apiVersion: v1
            kind: Pod
            name: my-pod
        # ...
        finally:
        - describe:
            apiVersion: v1
            kind: Pod
            name: my-pod
        # ...
    ```

If a `namespace` is specified, Chainsaw will describe resources in the specified namespace.

!!! example "Describe pods in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        catch:
        - describe:
            apiVersion: v1
            kind: Pod
            namespace: foo
        # ...
        finally:
        - describe:
            apiVersion: v1
            kind: Pod
            namespace: foo
        # ...
    ```

### Label selector

An optional [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) can be configured to refine the resources to be described.

!!! example "Describe pods using a label selector in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        catch:
        - describe:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
        # ...
        finally:
        - describe:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
        # ...
    ```

If a `namespace` is specified, Chainsaw will describe resources in the specified namespace.

!!! example "Describe pods using a label selector in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        catch:
        - describe:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
            namespace: foo
        # ...
        finally:
        - describe:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
            namespace: foo
        # ...
    ```

### Show events

The `showEvents` field can be used to enable or disable showing events when describing resources.

!!! note "Default"
    By default, `showEvents`is `true`.

!!! example "Do not show events"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        catch:
        - describe:
            apiVersion: v1
            kind: Pod
            namespace: foo
            showEvents: false
        # ...
        finally:
        - describe:
            apiVersion: v1
            kind: Pod
            namespace: foo
            showEvents: false
        # ...
    ```
