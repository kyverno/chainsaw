# Label selectors

Chainsaw can filter the tests to run using [label selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors).

You can pass label selectors using the `--selector` flag when invoking the `chainsaw test` command.

## Example

Given the test below:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: basic
  labels:
    foo: bar
spec:
  # ...
```

Invoking Chainsaw with the command below will take the test above into account:

```bash
$ chainsaw test --selector foo=bar
```
