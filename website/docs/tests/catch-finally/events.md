# Events

Collecting namespace events can help understand what happened inside the cluster.

## Configuration

The full structure of the `Events` resource is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Events).

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
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events:
            name: my-event
    ```

If a `namespace` is specified, Chainsaw will retrieve the specified event in the specified namespace.

!!! example "Collect event in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events:
            name: my-event
            namespace: foo
    ```

### All events

If no `name` and `namespace` is specified, Chainsaw will retrieve all events in the test namespace.

!!! example "Collect all events in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events: {}
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
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events:
            namespace: foo
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
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events:
            selector: app=my-app
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
        - apply:
            file: my-pod.yaml
        - assert:
            file: my-pod-assert.yaml
        catch:
        - events:
            selector: app=my-app
            namespace: foo
    ```
