# Command line flags

Even after a [configuration file](./file.md) is loaded, you can override specific settings using command-line flags.

## Example

```bash
chainsaw test                         \
  path/to/test/dir                    \
  --config path/to/your/config.yaml   \
  --assert-timeout 45s                \
  --skip-delete false                 \
  --fail-fast true                    \
  --parallel 4                        \
  ...
```

In this example, Chainsaw will load a configuration file but the timeout configuration and other settings will be overridden by the values set in the flags, regardless of the value in the loaded configuration file.

## Usage

The command below will run tests using the configuration from `my-config.yaml`, taking tests from `/path/to/tests`, and running a maximum of `10` tests simultaneously.

```bash
chainsaw test /path/to/tests --config my-config.yaml --parallel 10
```

## Reference documentation

See [Chainsaw test command reference](../commands/chainsaw_test.md#options) for more details.
