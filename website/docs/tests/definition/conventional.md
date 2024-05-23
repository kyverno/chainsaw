# Conventional approach

!!! warning
    While Chainsaw supports the conventional approach, we strongly recommend the explicit one.

    If you are new to Chainsaw we suggest you skip this section and jump directly to the [Explicit approach](./explicit.md).

## Introduction

The conventional approach is the simplest and less verbose one.

You provide bare Kubernetes resource manifests and Chainsaw will use those manifests to create, update, or assert expectations against a cluster.

### Limitations

While this syntax is simple, it suffers lots of limitations. It doesn't support deletion operations, commands, scripts, and all Chainsaw helpers.

It is also impossible to specify additional configuration per test, step or individual operation (timeouts, additional verifications, etc...), making this approach highly limited.

It also relies a lot on file naming conventions which can be error prone.

Finally, this approach doesn't encourage reusing files across tests and leads to duplication, making maintenance harder.

## File naming convention

Manifest files must follow a specific naming convention:
```
<step index>-<name|assert|errors>.yaml
```

As an example, `00-configmap.yaml`, `01-assert.yaml` and `02-errors.yaml` are valid file names.

### Assembling steps

It's perfectly valid to have multiple files for the same step.

Let's say we have the following files `00-resources.yaml`, `00-more-resources.yaml`, `00-assert.yaml` and `00-errors.yaml`:

- `00-resources.yaml` and `00-more-resources.yaml` contain resources that will be applied in step `00`
- `00-assert.yaml` contains assert statements in step `00`
- `00-errors.yaml` contains error statements in step `00`

With the four files above, Chainsaw will assemble a test step made of the combination of all those files.

### Loading process

The logic to determine the content of a step is always:

- The step index is obtained from the beginning of the file name, it must be composed of two numbers between `0` and `9` (from `00` to `99`)
- The next character acts as a separator and is expected to be `-`
- The rest of the file name (without extension) is then evaluated
    - If it is equal to `assert`, the content is considered assertion statements
    - If it is equal to `error`, the content is considered error statements
    - Else the content is considered resources to be applied
- The extension must be `.yaml` or `.yml`

## Example

### 01-configmap.yaml

The manifest below contains a config map in a file called `01-configmap.yaml`.
Chainsaw will associate this manifest with an apply operation in step `01`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
```

### 02-assert.yaml

The manifest below contains an assertion statement in a file called `02-assert.yaml`.
Chainsaw will associate this manifest with an assert operation in step `02`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
```

### 03-errors.yaml

The manifest below contains an error statement in a file called `03-errors.yaml`.
Chainsaw will associate this manifest with an error operation in step `03`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  lorem: ipsum
```

## Conclusion

This test will first create a config map, then assert the content of the config map contains the `foo: bar` data, and then verify that the config map does not contain the `lorem: ipsum` data.

For such a simple test, the conventional approach works reasonably well but will quickly become limited when the test scenarios get more complex.

Look at the [explicit approach](./explicit.md) for a lot more flexible solution.
