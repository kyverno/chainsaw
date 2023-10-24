# Chainsaw Tool Configuration

Chainsaw is a comprehensive tool designed to facilitate end-to-end (e2e) testing in Kubernetes. This documentation will focus on providing a breakdown of its configuration structure and how to use it.

Chainsaw can be configured in two ways:

- Using a configuration file
- Overriding configuration with command-line flags

If both are specified, **command-line flags will take precedence** over configuration coming from a configuration file.

## Configuration file

### Overview

Chainsaw is described as a "Stronger tool for e2e testing". With its versatile configuration options, you can customize the testing process to fit your needs.

### Configuration loading process

Chainsaw prioritizes its configuration in the following order:

1. **User-specified configuration**: If you explicitly provide a configuration file using a command-line flag.
1. **Default configuration file**: If no configuration is specified, Chainsaw will look for a default file named `.chainsaw.yaml` in the current working directory.
1. **Internal default configuration**: In the absence of both the above, Chainsaw will use a default configuration file embedded in the Chainsaw binary.

### How to specify a configuration

To use a custom configuration file:

```bash
chainsaw test --config path/to/your/config.yaml
```

If you don't specify any configuration, Chainsaw will look for the default configuration file `.chainsaw.yaml` in the current working directory. If that's not found, it will fall back to its internal default configuration.


### Example

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: custom-config
spec:
  timeout: 45s
  skipDelete: false
  failFast: true
  parallel: 4
  // ....
```

The full structure of the configuration file is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Configuration)

### Default configuration

The default configuration below is used by Chainsaw when no configuration file was provided and the default file `.chainsaw.yaml` does not exist.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: default
spec: {}
```

## Overriding with flags

Even after a configuration is loaded, you can override specific settings using command-line flags:

```bash
chainsaw test                           \
    --config path/to/your/config.yaml   \
    --test-dir path/to/test/dir         \
    --timeout 45s                       \
    --skip-delete false                 \
    --fail-fast true                    \
    --parallel 4                        \
    ...
```

In this example, Chainsaw will load a configuration file but the timeout configuration and other settings will be overridden by the values set in the flags, regardless of the value in the loaded configuration file.

Supported command line flags are documented [here](../commands/chainsaw_test.md#options)

## Usage example

```bash
chainsaw test --config my-config.yaml --test-dir /path/to/tests --parallel 10
```

This command will run tests using the configuration from `my-config.yaml`, taking tests from `/path/to/tests`, and running a maximum of `10` tests simultaneously.

## Reference

Refer to the reference documentations for details about supported fields in the configuration file and/or supported flags in the `test` command.

- [Configuration API reference](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Configuration)
- [Chainsaw test command reference](../commands/chainsaw_test.md#options)
