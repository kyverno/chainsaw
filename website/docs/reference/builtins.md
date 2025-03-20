# Built-in bindings

Chainsaw provides built-in bindings listed below.

## Common

| Name | Purpose | Type |
|---|---|---|
| `$values` | Values provided when invoking chainsaw with `--values` flag | `any` |
| `$namespace` | Name of the current test namespace | `string` |
| `$client` | Kubernetes client chainsaw is connected to (if not running with `--no-cluster`) | `object` |
| `$config` | Kubernetes client config chainsaw is connected to (if not running with `--no-cluster`) | `object` |

## In tests

| Name | Purpose | Type |
|---|---|---|
| `$test.id` | Current test id | `int` |
| `$test.metadata` | Current test metadata | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |

!!! note
    - `$test.id` starts at 1 for the first test

## In steps

| Name | Purpose | Type |
|---|---|---|
| `$step.id` | Current step id | `int` |

!!! note
    - `$step.id` starts at 1 for the first step

## In operations

| Name | Purpose | Type |
|---|---|---|
| `$operation.id` | Current operation id | `int` |
| `$operation.resourceId` | Current resource id | `int` |

!!! note
    - `$operation.id` starts at 1 for the first operation
    - `$operation.resourceId` maps to the resource id (starting at 1) in case the operation loads a file that contains multiple resources (the same operation is repeated once per resource)

## In checks and outputs

| Name | Purpose | Type |
|---|---|---|
| `@` | The state of the resource (if any) at the end of the operation | `any` |
| `$error` | The error message (if any) at the end of the operation | `string` |
| `$warnings` | List of warnings returned by the K8s API server for the operation | `{code: int, agent: string, text: string}` |
| `$stdout` | The content of the standard console output (if any) at the end of the operation | `string` |
| `$stderr` | The content of the standard console error output (if any) at the end of the operation | `string` |

!!! note
    - `$stdout` and `$stderr` are only available in `script` and `command` operations
    - `$warnings` are only available in `apply`, `create`, `patch`, `update`, and `delete` operations
