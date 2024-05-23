# Configuration file

Chainsaw prioritizes its configuration in the following order:

1. **User-specified configuration**

    If you explicitly provide a configuration file using a command-line flag

1. **Default configuration file**

    If no configuration is specified, Chainsaw will look for a default file named `.chainsaw.yaml` in the current working directory

1. **Internal default configuration**

    In the absence of both of the above, Chainsaw will use a [default configuration](#default-configuration) file embedded in the Chainsaw binary

## Example

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  timeouts:
    apply: 45s
    assert: 20s
    cleanup: 45s
    delete: 25s
    error: 10s
    exec: 45s
  cleanup:
    skipDelete: false
  execution:
    failFast: true
    parallel: 4
  # ...
```

## How to specify a configuration

To use a custom configuration file:

```bash
chainsaw test --config path/to/your/config.yaml
```

## Default configuration

The default configuration below is used by Chainsaw when no configuration file was provided and the default file `.chainsaw.yaml` does not exist.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: default
spec: {}
```

## Reference documentation

See [Configuration API reference](../reference/apis/chainsaw.v1alpha2.md#chainsaw-kyverno-io-v1alpha2-Configuration) for more details.
