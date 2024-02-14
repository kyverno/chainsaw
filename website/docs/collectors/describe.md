# Describe

Describing resources present in the cluster can help understand what happened and troubleshoot test failures.

## Configuration

The full structure of the `Describe` resource is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Describe).

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
            resource: pods
            name: my-pod
        # ...
        finally:
        - describe:
            resource: pods
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
            resource: pods
            namespace: foo
        # ...
        finally:
        - describe:
            resource: pods
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
            resource: pods
            selector: app=my-app
        # ...
        finally:
        - describe:
            resource: pods
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
            resource: pods
            selector: app=my-app
            namespace: foo
        # ...
        finally:
        - describe:
            resource: pods
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
            resource: pods
            namespace: foo
            showEvents: false
        # ...
        finally:
        - describe:
            resource: pods
            namespace: foo
            showEvents: false
        # ...
    ```
