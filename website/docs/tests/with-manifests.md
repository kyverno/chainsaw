# Manifests based syntax

This is the simplest and less verbose supported syntax, you provide bare Kubernetes resource manifests and Chainsaw will use those manifests to create, update, or assert expectations against a cluster.

While this syntax is simple, it doesn't support deletion operations and doesn't allow specifying additional configuration per test or step.

It also relies a lot on file naming conventions and makes it hard to reuse files across multiple tests.

## File naming convention

Manifest files must follow the naming convention `<step index>-<name|assert|error>.yaml`.

As an example `00-configmap.yaml`, `01-assert.yaml` and `02-error.yaml` can all be considered valid file names.

It's also perfectly valid to have multiple files for the same step. Let's imagine we have the following files `00-resources.yaml`, `00-more-resources.yaml`, `00-assert.yaml` and `00-error.yaml`:

- `00-resources.yaml` and `00-more-resources.yaml` contain resources that will be applied in step `00`
- `00-assert.yaml` contains assert statements in step `00`
- `00-error.yaml` contains error statements in step `00`

With the four files above, Chainsaw will assemble a test step made of the combination of all those files.

## Manifests loading process

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

The manifest below contains a config map in a file called `01-configmap.yaml`. Chainsaw will associate this manifest with an apply operation in step `01`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
```

### 02-assert.yaml

The manifest below contains an assertion statement in a file called `02-assert.yaml`. Chainsaw will associate this manifest with an assert operation in step `02`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
```

### 03-error.yaml

The manifest below contains an error statement in a file called `03-error.yaml`. Chainsaw will associate this manifest with an error operation in step `03`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  lorem: ipsum
```

### Conclusion

This test will first create a config map, then assert the content of the config map contains the `foo: bar` data, and then verify that the configmap does not contain the `lorem: ipsum` data.

Assert and error statements are very similar to standard resource definitions, they share the same structure but can be partial. Chainsaw will process only the elements present in the statements against the actual resource definition.