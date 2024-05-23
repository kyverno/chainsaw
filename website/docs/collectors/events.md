# Events

Collecting namespace events can help understand what happened inside the cluster.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Events` resource is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Events).

### All namespaces

When used with a namespaced resource, it is possible to consider all namespaces in the cluster by setting `namespace: '*'`.

## Usage examples

### Single event

If a `name` is specified, Chainsaw will retrieve the specified event in the test namespace.

!!! example "Collect event in the test namespace"

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
        - events:
            name: my-event
        # ...
        finally:
        - events:
            name: my-event
        # ...
    ```

If a `namespace` is specified, Chainsaw will retrieve the specified event in the specified namespace.

!!! example "Collect event in a specific namespace"

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
        - events:
            name: my-event
            namespace: foo
        # ...
        finally:
        - events:
            name: my-event
            namespace: foo
        # ...
    ```

### All events

If no `name` and `namespace` are specified, Chainsaw will retrieve all events in the test namespace.

!!! example "Collect all events in the test namespace"

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
        - events: {}
        # ...
        finally:
        - events: {}
        # ...
    ```

On the other hand, if a `namespace` is specified, Chainsaw will retrieve all events in the specified namespace.

!!! example "Collect all events in a specific namespace"

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
        - events:
            namespace: foo
        # ...
        finally:
        - events:
            namespace: foo
        # ...
    ```

### Label selector

An optional [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) can be configured to refine the events to be retrieved.

!!! example "Collect events using a label selector in the test namespace"

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
        - events:
            selector: app=my-app
        # ...
        finally:
        - events:
            selector: app=my-app
        # ...
    ```

If a `namespace` is specified, Chainsaw will retrieve events using the specified namespace.

!!! example "Collect events using a label selector in a specific namespace"

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
        - events:
            selector: app=my-app
            namespace: foo
        # ...
        finally:
        - events:
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
        - events:
            format: json
        # ...
        finally:
        - events:
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
        - events:
            format: yaml
        # ...
        finally:
        - events:
            format: yaml
        # ...
    ```
