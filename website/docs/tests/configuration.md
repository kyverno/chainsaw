# Configuration

Chainsaw can be configured in two ways:

- Using a configuration file
- Using command line flags

If both are specified, **command line flags will take precedence** over configuration coming from a configuration file.

## Using a configuration file

You can provide a configuration file when invoking Chainsaw with the `--config <path to config file>` flag.

The config file is structured as a Kubernetes manifest with:

- `apiversion`: `chainsaw.kyverno.io/v1alpha1`
- `kind`: `Configuration`

The structure of the configuration files is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Configuration)

Note that **Chainsaw will always load a configuration file** and use the logic below to determine the file to load:

1. If the `--config` flag is present, use the specified configuration file
1. If the `--config` flag is NOT present, try with the [default](#default-configuration) `.chainsaw.yaml` file name in the current working directory
1. Load a default configuration file embedded in the Chainsaw binary

It works this way because Chainsaw uses an advanced resource loading mechanism and supports most validation and defaulting markers.

**Example:**

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

### Default configuration

The default configuration below is used by Chainsaw when no configuration file was provided and the default file `.chainsaw.yaml` does not exist.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: default
spec: {}
```

## Using command line flags

All config elements can also be specified using command line flags.

You can provide a combination of flags when invoking Chainsaw, and even combine a configuration file with flags, in this case flags will take precedence over configuration elements coming from the config file.

**Example:**

```bash
chainsaw test                           \
    --test-dir ./path/to/test/folder    \
    --timeout 45s                       \
    --skip-delete false                 \
    --fail-fast true                    \
    --parallel 4                        \
    ...
```

Supported command line flags are documented [here](../commands/chainsaw_test.md#options)

## Per test (or step) options

The configuration options documented here are considered the global configuration.

Some of them can be overridden at the test or test step level and will be discussed in the next sections.