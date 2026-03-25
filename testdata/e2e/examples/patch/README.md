# Test: `patch`

*No description*

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 0 | 3 | 0 | 0 | 0 |
| 2 | [step-2](#step-step-2) | 0 | 3 | 0 | 0 | 0 |

### Step: `step-1`

Patch a ConfigMap

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 0 | 0 | *No description* |
| 2 | `patch` | 0 | 0 | *No description* |
| 3 | `assert` | 0 | 0 | *No description* |

### Step: `step-2`

Patch the 'status' subresource of a Pod

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 0 | 0 | *No description* |
| 2 | `patch` | 0 | 0 | *No description* |
| 3 | `assert` | 0 | 0 | *No description* |

---

