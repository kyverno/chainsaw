# Passing data to tests

Chainsaw can pass arbitrary values when running tests using the `--values` flag.
Values will be available to tests under the `$values` binding.

This is useful when a test needs to be configured externally.

## Example

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

Now you can invoke chainsaw like this:

```bash
# pass object { "foo": "bar" } as values to the executed tests
# `--values -` means values are read from standard input
echo "foo: bar" | chainsaw test --values -

# read values from a file
chainsaw test --values ./values.yaml

# pass values using heredoc
chainsaw test --values - <<EOF
foo: bar
EOF
```
