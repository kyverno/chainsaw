# Generating test docs

## Overview

Chainsaw makes it simple to generate the documentation of your tests.

As test suites grow, it becomes important to document what a test does and how it is supposed to work.

Going through the implementation of a test to understand its purpose is not an efficient strategy.

!!! tip "Reference documentation"

    You can view the full command documentation [here](../commands/chainsaw_generate_docs.md).

## Usage

To generate the docs of a test, Chainsaw provides the `chainsaw generate docs` command.

```bash
chainsaw generate docs --test-dir path/to/chainsaw/tests
```

This will automatically discover tests and document steps and operations in `try`, `catch` and `finally` statements.

## The `description` field


Additionally, you can set the `description` field in:

- `TestSpec`
- `TestStepSpec`
- `Operation`
- `Catch`
- `Finally`

Chainsaw will output them nicely in the generated docs.

## Example

See below for an example test and the corresponding generated docs.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: basic
spec:
  description: This is a very simple test that creates a configmap and checks the content is as expected.
  steps:
  - description: This steps applies the configmap in the cluster and checks the configmap content.
    try:
    - description: Create the configmap.
      apply:
        file: configmap.yaml
    - description: Check the configmap content.
      assert:
        file: configmap-assert.yaml
```

---

# Test: `basic`

This is a very simple test that creates a configmap and checks the content is as expected.

### Steps

| # | Name | Try | Catch | Finally |
|:-:|---|:-:|:-:|:-:|
| 1 | [step-1](#step-step-1) | 2 | 0 | 0 |

## Step: `step-1`

This steps applies the configmap in the cluster and checks the configmap content.

### Try

| # | Operation | Description |
|:-:|---|---|
| 1 | `apply` | Create the configmap. |
| 2 | `assert` | Check the configmap content. |
