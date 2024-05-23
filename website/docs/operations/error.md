# Error

The `error` operation lets you define a set of expected errors for a test step. If any of these errors occur during the test, they are treated as expected outcomes. However, if an error that's not on this list occurs, it will be treated as a test failure.

!!! info "Assertion trees"

    Errors in Chainsaw are based on **assertion trees**.

    Assertion trees is a solution to declaratively represent complex conditions like partial array comparisons or complex operations against an incoming data structure.

    Assertion trees are compatible with standard assertions that exist in tools like KUTTL but can do a lot more.
    Please see the [assertion trees documentation](https://kyverno.github.io/kyverno-json/latest/policies/asserts/) in kyverno-json for details.

## Configuration

!!! tip "Reference documentation"
    - The full structure of the `Error` is documented [here](../reference/apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Error).
    - This operation supports [bindings](../tests/general/bindings/index.md).

## Usage examples

Below is an example of using `error` in a `Test` resource.

!!! example "Expecting an error from a specific file"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - error:
            file: ../resources/deployment-error.yaml
        # ...
    ```

!!! example "Using file path expressions to expect errors from multiple files"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example-multi
    spec:
      steps:
      - try:
        # ...
        - error:
            file: "../errors/*.yaml"
        # ...
    ```

!!! example "Using an URL"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - error:
            file: https://raw.githubusercontent.com/user/repo/branch/path/to/deployment-error.yaml
        # ...
    ```

!!! example "Using an inline assertion tree"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: example
    spec:
      steps:
      - try:
        # ...
        - error:
            resource:
              apiVersion: v1
              kind: Deployment
              metadata:
                name: foo
              spec:
                (replicas > 3): true
        # ...
    ```
