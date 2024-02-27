# Wait

The `wait` is a wrapper around `kubectl wait`. It allows to wait for deletion or conditions against resources.

!!! tip "Reference documentation"
    The full structure of the `Wait` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Wait).

## Usage in `Test`

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
            resource: pods
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
            resource: pods
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
            resource: pods
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
            resource: pods
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
            resource: pods
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
            resource: pods
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
            resource: pods
            format: json
            # ...
        catch:
        # ...
        - wait:
            resource: pods
            format: json
            # ...
        finally:
        # ...
        - wait:
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
        - wait:
            resource: pods
            format: yaml
            # ...
        catch:
        # ...
        - wait:
            resource: pods
            format: yaml
            # ...
        finally:
        # ...
        - wait:
            resource: pods
            format: yaml
            # ...
    ```
