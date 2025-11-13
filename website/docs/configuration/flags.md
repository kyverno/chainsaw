# Command line flags

After a [configuration file](./file.md) is loaded, you can override specific settings using command-line flags.

!!! info "Precedence"
    Command-line flags always take precedence over the configuration coming from a configuration file.

## Example

```bash
chainsaw test                         \
  path/to/test/dir                    \
  --config path/to/your/config.yaml   \
  --assert-timeout 45s                \
  --skip-delete false                 \
  --fail-fast true                    \
  --parallel 4                        \
  --values ./values.yaml              \
  --set env=production                \
  --set-string image.tag=v1.2.3       \
  ...
```

In this example, Chainsaw will load a configuration file but the timeout configuration and other settings will be overridden by the values set in the flags, regardless of the value in the loaded configuration file.

## Reference documentation

See [chainsaw test](../reference/commands/chainsaw_test.md#options) command reference for the list of all available flags.
