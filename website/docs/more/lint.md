# Lint tests

## Overview

Chainsaw comes with a `lint` command to detect ill-formated tests.

!!! tip "Reference documentation"

    You can view the full command documentation [here](../commands/chainsaw_lint.md).

## Usage

To build the docs of a test, Chainsaw provides the `chainsaw lint test -f path/to/chainsaw-test.yaml` command.

```bash
chainsaw lint test -f - <<EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: assertion-tree
spec:
  steps:
  - try:
    - assert:
        file: assert.yaml
EOF
```

```bash
Processing input...
The document is valid
```
