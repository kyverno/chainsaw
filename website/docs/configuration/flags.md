# Command line flags

After a [configuration file](./file.md) is loaded, you can override specific settings using command-line flags.

## Reference documentation

See [Chainsaw test command reference](../commands/chainsaw_test.md#options) for the list of all available flags.

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

## Configuration in logs

Chainsaw will print its configuration at startup:

```
Version: (devel)
Loading default configuration...
- Using test file: chainsaw-test
- TestDirs [./testdata/e2e/examples/dynamic-clusters]
- SkipDelete false
- FailFast false
- ReportFormat ''
- ReportName ''
- Namespace ''
- FullName false
- IncludeTestRegex ''
- ExcludeTestRegex ''
- ApplyTimeout 5s
- AssertTimeout 30s
- CleanupTimeout 30s
- DeleteTimeout 15s
- ErrorTimeout 30s
- ExecTimeout 5s
- NoCluster false
- PauseOnFailure false
...
```