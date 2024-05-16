# Resource templating

!!! warning "Experimental status"
    This is an **experimental feature**, and implementation could change slightly in the next versions.

!!! info
    Templating was disabled by default in `v0.1.*` but is now enabled by default since `v0.2.1`.

Chainsaw can apply transformations to the resources before they are processed by the operation.

This is useful when a resource needs some runtime configuration.

Templating must be enabled at the configuration, test, step, or operation level for the templating process to kick in.
Alternatively, templating can be enabled using the `--template` flag when invoking chainsaw from the command line.

!!! note
    Unlike assertion trees, templating can only be specified in leave nodes of the YAML tree.

## Supported operations

Resource templating is supported in the following operations:

- [Apply](./apply.md)
- [Assert](./assert.md)
- [Create](./create.md)
- [Delete](./delete.md)
- [Error](./error.md)
- [Patch](./patch.md)
- [Update](./update.md)

### Assert and Error

When templating `assert` or `error` operations, the content is already an assertion tree.

For this reason, only the elements used for looking up the resources to be processed by the operation will be considered for templating.
That is, only `apiVersion`, `kind`, `name`, `namespace` and `labels` are considered for templating.
Other fields are not, they are part of the assertion tree only.

!!! example "assert and error example"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: template
    spec:
      template: true
      steps:
      - assert:
          resource:
            # apiVersion, kind, name, namespace and labels are considered for templating
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: ($namespace)
            # other fields are not (they are part of the assertion tree)
            data:
              foo: ($namespace)
      - error:
          resource:
            # apiVersion, kind, name, namespace and labels are considered for templating
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: ($namespace)
            # other fields are not (they are part of the assertion tree)
            data:
              bar: ($namespace)
    ```

### Apply, Create and Patch

When templating `apply`, `create` or `patch` operations, the whole content is considered for templating.

!!! example "apply and create example"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: template
    spec:
      template: true
      steps:
      - apply:
          resource:
            # the whole content is considered for templating
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: ($namespace)
            data:
              foo: ($namespace)
      - create:
          resource:
            # the whole content is considered for templating
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: ($namespace)
            data:
              foo: ($namespace)
      - patch:
          resource:
            # the whole content is considered for templating
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: ($namespace)
            data:
              foo: ($namespace)
    ```

### Delete

When templating `delete` operations, the whole content is considered for templating.

!!! example "apply and create example"

    ```yaml
    apiVersion: chainsaw.kyverno.io/v1alpha1
    kind: Test
    metadata:
      name: template
    spec:
      template: true
      steps:
      - delete:
          ref:
            # the whole content is considered for templating
            apiVersion: v1
            kind: ConfigMap
            name: ($namespace)
    ```