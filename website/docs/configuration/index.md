# Configuring Chainsaw

This documentation focuses on providing a breakdown of the Chainsaw configuration structure and how to use it.

Chainsaw can be configured in two different and complementary ways:

- [Using a configuration file](./file.md)
- [Overriding configuration with command-line flags](./flags.md)

!!! tip "Precedence"
    If both are specified, **command-line flags will take precedence** over the configuration coming from a configuration file.

## Specific configuration options

Please pay attention to the configuration options below, they may or may not be relevant in your case but can be useful in certain cases:

- [Timeouts](./options/timeouts.md)
- [Discovery options](./options/discovery.md)
- [Execution options](./options/execution.md)
- [Namespace options](./options/namespace.md)
- [Templating options](./options/templating.md)
- [Cleanup options](./options/cleanup.md)
- [Error options](./options/error.md)
- [Reporting options](./options/report.md)
- [Multi cluster options](./options/clusters.md)

TODO

- [Termination graceful period](./options/grace.md)
- [Passing arbitrary values to tests](./options/values.md)
- [Running without a cluster](./options/no-cluster.md)
