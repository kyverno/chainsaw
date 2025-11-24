# External values

Chainsaw can pass arbitrary values when running tests using the `--set`, `--set-string` and `--values` flag.
Values will be available to tests under the `$values` binding.

## Configuration

### With file

!!! note
    Values can't be configured with a configuration file.

### With flags

```bash
chainsaw test --values ./values.yaml \
    --set env=poc \
    --set clusterDirectory=my-cluster \
    --set-string image.tag=01
```

#### Command line overrides

- `--set key=value` — set or override a value (parses types)
- `--set-string key=value` — set or override a value as string

```bash
chainsaw test --values ./values.yaml --set env=poc --set-string image.tag=01
```
