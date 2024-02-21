# Get

The `get` collector is used to list and print resources in the cluster.

## Configuration

The full structure of the `Get` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Get).

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
            resource: pods
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
            resource: pods
            name: my-pod
            namespace: foo
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
            namespace: foo
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
            selector: app=my-app
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
            selector: app=my-app
            namespace: foo
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
            format: json
        # ...
        finally:
        - get:
            resource: pods
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
            resource: pods
            format: yaml
        # ...
        finally:
        - get:
            resource: pods
            format: yaml
        # ...
    ```
