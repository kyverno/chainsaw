# Configuring Chainsaw

This documentation focuses on providing a breakdown of the Chainsaw configuration structure and how to use it.

Chainsaw can be configured in two different and complementary ways:

- [Using a configuration file](./file.md)
- [Overriding configuration with command-line flags](./flags.md)

!!! tip "Precedence"
    If both are specified, **command-line flags will take precedence** over configuration coming from a configuration file.

## Specific configuration options

Please pay attention to the configuration options below, they may or may not be relevant in your case but can be useful in certain cases:

- [Timeouts](./options/timeouts.md)
- [Termination graceful period](./options/grace.md)
- [Delay before cleanup](./options/cleanup-delay.md)
- [Creating test reports](./options/reports.md)
- [Label selectors](./options/selector.md)
- [Passing arbitrary values to tests](./options/values.md)
- [Multi cluster](./options/multi-cluster.md)
- [Resource templating](./options/templating.md)
- [Running without a cluster](./options/no-cluster.md)
