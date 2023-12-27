# Configuring Chainsaw

Chainsaw is a comprehensive tool designed to facilitate **end-to-end testing in [Kubernetes](https://kubernetes.io)**.

This documentation will focus on providing a breakdown of its configuration structure and how to use it.

Chainsaw can be configured in two different and complementary ways:

- [Using a configuration file](./file.md)
- [Overriding configuration with command-line flags](./flags.md)

!!! note "Precedence"
    If both are specified, **command-line flags will take precedence** over configuration coming from a configuration file.

## Specific configuration options

Please pay attention to the configuration options below, they may or may not be relevant in your case but can be useful in certain cases:

- [Timeouts](./timeouts.md)
- [Termination graceful period](./grace.md)
- [Delay before cleanup](./cleanup-delay.md)
- [Label selectors](./selector.md)
