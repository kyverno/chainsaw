# No cluster options

Chainsaw can be run without any connection to a Kubernetes cluster.

In this case, Chainsaw will not try to create an ephemeral namespace and all operations requiring a Kubernetes cluster will fail.

## Configuration

### With file

!!! note
    No cluster options can't be configured with a configuration file.

### With flags

```bash
chainsaw test --no-cluster
```
