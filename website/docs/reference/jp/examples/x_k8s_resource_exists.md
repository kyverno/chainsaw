# x_k8s_resource_exists

## Signature

`x_k8s_resource_exists(any, string, string)`

## Description

Checks if a given resource type is available in a Kubernetes cluster.

## Examples

!!! info "Clustered resources"

    For clustered resources, you can leave the namespace empty `''`.

```
# `$client` is a binding pointing to a Kubernetes client

x_k8s_resource_exists($client, 'apps/v1', 'Deployment')
```
