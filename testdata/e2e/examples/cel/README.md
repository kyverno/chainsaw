# Test: `cel`

*No description*

## Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `a` | 1 |

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 1 | 1 | 0 | 0 | 0 |
| 2 | [step-2](#step-step-2) | 1 | 1 | 0 | 0 | 0 |

### Step: `step-1`

*No description*

#### Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `b` | 2 |

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 1 | 0 | *No description* |

### Step: `step-2`

*No description*

#### Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `b` | 2 |

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `assert` | 1 | 0 | *No description* |

---

