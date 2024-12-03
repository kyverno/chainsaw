# Test: `basic`

This is a very simple test that creates a configmap and checks the content is as expected.

## Steps

| # | Name | Bindings | Try | Catch | Finally | Cleanup |
|:-:|---|:-:|:-:|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 0 | 3 | 0 | 0 | 0 |

### Step: `step-1`

This steps applies the configmap in the cluster and checks the configmap content.

#### Try

| # | Operation | Bindings | Outputs | Description |
|:-:|---|:-:|:-:|---|
| 1 | `apply` | 0 | 0 | Create the configmap. |
| 2 | `assert` | 0 | 0 | Check the configmap content. |
| 3 | `script` | 0 | 0 | *No description* |

---

