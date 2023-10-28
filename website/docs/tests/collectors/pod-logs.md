# Pod logs

Collecting pod logs can help understand what happened inside one or more pods.

## Configuration

The full structure of the `PodLogs` resource is documented [here](../../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-PodLogs).

TODO `Container` and `Tail`

### Single pod

If a pod `name` was specified, Chainsaw will retrieve logs from this specific pod in the test namespace.

!!! note "Collect pod logs in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs:
                name: my-pod
    ```

If a `namespace` was specified, Chainsaw will retrieve logs from this specific pod in the specified namespace.

!!! note "Collect pod logs in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs:
                name: my-pod
                namespace: foo
    ```

### All pods

If no pod `name` and `namespace` was specified, Chainsaw will retrieve logs from all pods in the test namespace.

!!! note "Collect all pod logs in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs: {}
    ```

On the other hand, if a `namespace` was specified, Chainsaw will retrieve logs from all pods in the specified namespace.

!!! note "Collect all pod logs in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs:
                namespace: foo
    ```

### Label selector

An optional [label selector](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors) can be configured to refine the pods to retrieve logs from.

!!! note "Collect pod logs using a label selector in the test namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs:
                selector: app=my-app
    ```

If a `namespace` was specified, Chainsaw will retrieve pod logs using the specified namespace.

!!! note "Collect pod logs using a label selector in a specific namespace"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - spec:
          apply:
          - file: my-pod.yaml
          assert:
          - file: my-pod-assert.yaml
          onFailure:
          - collect:
              podLogs:
                selector: app=my-app
                namespace: foo
    ```
