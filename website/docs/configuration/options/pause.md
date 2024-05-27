# Pause options

Chainsaw can be configured to pause and wait for user input when a failure happens.
This is useful when Chainsaw is run locally to allow debugging and troubleshooting failures.

## Configuration

### With file

!!! note
    Pause options can't be configured with a configuration file.

### With flags

```bash
chainsaw test --pause-on-failure
```
