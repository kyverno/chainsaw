# External values

Chainsaw can pass arbitrary values when running tests using the `--values` flag.
Values will be available to tests under the `$values` binding.

## Configuration

### With file

!!! note
    Values can't be configured with a configuration file.

### With flags

```bash
chainsaw test --values ./values.yaml
```
