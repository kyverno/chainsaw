# Operation checks

Considering an operation's success or failure is not always as simple as checking an error code.

- Sometimes an operation can fail but the failure is what you expected, hence the operation should be reported as successful.
- Sometimes an operation can succeed but the result is not what you expected, in this case, the operation should be reported as a failure.

To support those kinds of use cases, some operations support an additional `check` field to evaluate the operation result against an assertion tree.

!!! info "Assertion trees"

    Assertions in Chainsaw are based on **assertion trees**.

    Assertion trees are a solution to declaratively represent complex conditions like partial array comparisons or complex operations against an incoming data structure.

    Assertion trees are compatible with standard assertions that exist in tools like KUTTL but can do a lot more.
    Please see the [assertion trees documentation](https://kyverno.github.io/kyverno-json/policies/asserts/) in kyverno-json for details.

!!! tip "Checked model"
    Different operations have a different model passed through the assertion tree.

    The object passed to the assertion tree is the output object of the operation. Additional data like error or standard logs are passed using bindings (`$error`, `$stdout`, `$stderr`)

## `Expect` vs `Check`

While a simple check is enough to determine the result of a single operation, we needed a more advanced construct to cover `apply` and `create` operations. Those operations can operate on files containing multiple manifests and every manifest can have a different result.

To support more granular checks we use the `expect` field that contains an array of [Expectations](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Expectation).
Every expectation is made of an optional `match` and a `check` statement.

This way it is possible to control the scope of a `check`.

!!! tip "Null match"
    If the `match` statement is null, the `check` statement applies to all manifests in the operation.

    If no expectation matches a given manifest, the default expectation will be used, checking that no error occurred.

## Apply

`apply` supports `expect` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `@` | The state of the resource (if any) at the end of the operation | `object` |

## Command

`command` supports `check` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `$stdout` | The content of the standard console output (if any) at the end of the operation | `string` |
| `$stderr` | The content of the standard console error output (if any) at the end of the operation | `string` |
| `@` | Always `null` | |

## Create

`create` supports `expect` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `@` | The state of the resource (if any) at the end of the operation | `object` |

## Delete

`delete` supports `check` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `@` | The state of the resource (if any) at the end of the operation | `object` |

## Patch

`patch` supports `expect` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `@` | The state of the resource (if any) at the end of the operation | `object` |

## Script

`script` supports `check` and has the following elements to be checked:

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `object` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `$stdout` | The content of the standard console output (if any) at the end of the operation | `string` |
| `$stderr` | The content of the standard console error output (if any) at the end of the operation | `string` |
| `@` | Always `null` | |
