# Label selectors

Chainsaw can filter the tests to run using [label selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors).

## Configuration

### With file

!!! note
    Label selectors can't be configured with a configuration file.

### With flags

```bash
chainsaw test --selector foo=bar
```
