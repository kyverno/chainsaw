# Test: `metrics-decode`

*No description*

## Bindings

| # | Name | Value |
|:-:|---|---|
| 1 | `metrics` | "# Only a quite simple scenario with two metric families.\n# More complicated tests of the parser itself can be found in the text package.\n# TYPE mf2 counter\nmf2 3\nmf1{label=\"value1\"} -3.14 123456\nmf1{label=\"value2\"} 42\nmf2 4\n" |

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 0 | 2 | 0 | 0 | 0 |

### Step: `step-1`

*No description*

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `script` | 0 | 0 | *No description* |
| 2 | `assert` | 0 | 0 | *No description* |

---

