# Test: `basic`

This is a very simple test that creates a configmap and checks the content is as expected.

### Steps

| # | Name | Try | Catch | Finally |
|:-:|---|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 3 | 0 | 0 |

## Step: `step-1`

This steps applies the configmap in the cluster and checks the configmap content.

### Try

| # | Operation | Description |
|:-:|---|---|
| 1 | `apply` | Create the configmap. |
| 2 | `assert` | Check the configmap content. |
| 3 | `script` | *No description* |
