# Wait

The `wait` operation is a wrapper around `kubectl wait`. It allows to wait for deletion or conditions against resources.

## Configuration

!!! tip "Reference documentation"
    The full structure of the `Wait` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Wait).

!!! warning "Deprecated syntax"
    You can specify the `resource` directly instead of using `apiVersion` and `kind`.
    
    **This is a deprecated syntax though and will be removed in a future version.**

### Clustered resources

When used with a clustered resource, the `namespace` is ignored and is not added to the corresponding `kubectl` command.

### All resources

If you don't specify a `name` or a `selector`, the `wait` operation will consider `all` resources.

### All namespaces

When used with a namespaced resource, it is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

## Usage examples

Below is an example of using `wait` in a `Test` resource.

!!! example "Wait pod ready"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            name: my-pod
            timeout: 1m
            for:
              condition:
                name: Ready
                value: 'true'
        # ...
    ```

!!! example "Wait pod ready in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            name: my-pod
            namespace: my-ns
            timeout: 1m
            for:
              condition:
                name: Ready
                value: 'true'
        # ...
    ```

!!! example "Wait pods ready using a label selector"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            selector: app=foo
            timeout: 1m
            for:
              condition:
                name: Ready
                value: 'true'
        # ...
    ```

!!! example "Wait pod deleted"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            name: my-pod
            timeout: 1m
            for:
              deletion: {}
        # ...
    ```

!!! example "Wait pod deleted in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            name: my-pod
            namespace: my-ns
            timeout: 1m
            for:
              deletion: {}
        # ...
    ```

!!! example "Wait pods deleted using a label selector"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            selector: app=foo
            timeout: 1m
            for:
              deletion: {}
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
        - wait:
            apiVersion: v1
            kind: Pod
            format: json
            # ...
        catch:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            format: json
            # ...
        finally:
        # ...
        - wait:
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
        - wait:
            apiVersion: v1
            kind: Pod
            format: yaml
            # ...
        catch:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            format: yaml
            # ...
        finally:
        # ...
        - wait:
            apiVersion: v1
            kind: Pod
            format: yaml
            # ...
    ```
