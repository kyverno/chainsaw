# Pass data to tests

Chainsaw can pass arbitrary values when running tests using the `--values` flag.
Values will be available to tests under the `$values` binding.

This is useful when a test needs to be configured externally.

## Reference external data

The test below expects the `$value.foo` to be provided when chainsaw is invoked.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - assert:
        resource:
          ($values.foo): bar
```

## Invoking Chainsaw

### Read values from a file

```bash
chainsaw test --values ./values.yaml
```

### Set values from command line

```bash
chainsaw test --set env=production --set-string version=v1.2.0
```

### Read from stdin

```bash
echo "foo: bar" | chainsaw test --values -
```

### Use heredoc

```bash
chainsaw test --values - <<EOF
foo: bar
EOF
```
