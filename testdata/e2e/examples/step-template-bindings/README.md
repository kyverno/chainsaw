# Test: `step-template-bindings`

*No description*

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 2 | 1 | 0 | 0 | 0 |
| 2 | [step-2](#step-step-2) | 0 | 1 | 0 | 0 | 0 |

### Step: `step-1`

*No description*

#### Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `input` | "from-template" |
| 2 | `input` | "from-test" |

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `create` | 0 | 0 | *No description* |

### Step: `step-2`

*No description*

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `assert` | 0 | 0 | *No description* |

---

