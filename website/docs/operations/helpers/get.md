# Get

The `get` collector is used to list and print resources in the cluster.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Get` resource is documented [here](../../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Get).

!!! warning "Deprecated syntax"
    You can specify the `resource` directly instead of using `apiVersion` and `kind`.
    
    **This is a deprecated syntax though and will be removed in a future version.**

### Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

### All namespaces

When used with a namespaced resource, it is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

## Usage examples

### Single resource

If a `name` is specified, Chainsaw will retrieve the specified resource in the test namespace.

!!! example "Get pod in the test namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
            name: my-pod
        # ...
        finally:
        - get:
            resource: pods
            name: my-pod
        # ...
    ```

If a `namespace` is specified, Chainsaw will retrieve the specified resource in the specified namespace.

!!! example "Collect pod in a specific namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
            name: my-pod
            namespace: foo
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            name: my-pod
            namespace: foo
        # ...
    ```

### All resources

If no `name` and `namespace` are specified, Chainsaw will retrieve all resources in the test namespace.

!!! example "Collect all resources in the test namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
        # ...
    ```

On the other hand, if a `namespace` is specified, Chainsaw will retrieve all resources in the specified namespace.

!!! example "Collect all resources in a specific namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
            namespace: foo
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            namespace: foo
        # ...
    ```

### Label selector

An optional [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) can be configured to refine the resources to be retrieved.

!!! example "Collect resources using a label selector in the test namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
        # ...
    ```

If a `namespace` is specified, Chainsaw will retrieve resources using the specified namespace.

!!! example "Collect resources using a label selector in a specific namespace"

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
        - get:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
            namespace: foo
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            selector: app=my-app
            namespace: foo
        # ...
    ```

### Format

An optional `format` can be specified. Supported formats are `json` and `yaml`.

If `format` is not specified, results will be returned in text format.

!!! example "Use json format"

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
        - get:
            apiVersion: v1
            kind: Pod
            format: json
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            format: json
        # ...
    ```

!!! example "Use yaml format"

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
        - get:
            apiVersion: v1
            kind: Pod
            format: yaml
        # ...
        finally:
        - get:
            apiVersion: v1
            kind: Pod
            format: yaml
        # ...
    ```
