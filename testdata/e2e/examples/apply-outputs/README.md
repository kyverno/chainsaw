# Test: `apply-outputs-1`

*No description*

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 0 | 3 | 0 | 0 | 0 |

### Step: `step-1`

*No description*

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 0 | 1 | *No description* |
| 2 | `script` | 0 | 0 | *No description* |
| 3 | `assert` | 0 | 0 | *No description* |

---

# Test: `apply-outputs-2`

*No description*

## Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `prefix` | "('prefix')" |
| 2 | `($prefix)` | "foo" |

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 1 | 6 | 0 | 0 | 0 |

### Step: `step-1`

*No description*

#### Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `foos` | [] |

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 0 | 2 | *No description* |
| 2 | `script` | 0 | 0 | *No description* |
| 3 | `script` | 0 | 0 | *No description* |
| 4 | `script` | 0 | 0 | *No description* |
| 5 | `assert` | 0 | 0 | *No description* |
| 6 | `assert` | 0 | 0 | *No description* |

---

