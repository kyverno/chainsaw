# Running without a cluster

Chainsaw can be run without connection to a Kubernetes cluster.
In this case, chainsaw will not try to create an ephemeral namespace and all operations requiring a Kubernetes cluster will fail.

To run chainsaw in this mode pass the `--no-cluster` flag.

## Example

```bash
# run chainsaw without connection to a Kubernetes cluster
chainsaw test --no-cluster
```
