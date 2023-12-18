# Setup

To use Chainsaw you will need a [Kubernetes](https://kybernetes.io) cluster, Chainsaw won't create one for you.

In these examples, we will use [kind](https://kind.sigs.k8s.io) but feel free to use the tool of your choice.

!!! warning "Not a cluster management tool"

    We consider this is not the responsibility of Chainsaw to manage clusters.

    There are plenty of solutions to create and manage local clusters that will do that better than Chainsaw.

## Create a [kind](https://kind.sigs.k8s.io) cluster

Please refer to the [kind install docs](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) to install it locally.

Once [kind](https://kind.sigs.k8s.io) creating a local cluster is as simple as running:

```bash
# define kind image
export KIND_IMAGE="kindest/node:v1.28.0"

# create cluster
kind create cluster --image $KIND_IMAGE
```

## Install Chainsaw

Refer to [install docs](../install.md) to install Chainsaw.
