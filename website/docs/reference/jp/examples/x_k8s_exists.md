# x_k8s_exists

## Signature

`x_k8s_exists(any, string, string, string, string)`

## Description

Checks if a given resource exists in a Kubernetes cluster.

## Examples

!!! info "Clustered resources"

    For clustered resources, you can leave the namespace empty `''`.

```
# `$client` is a binding pointing to a Kubernetes client

x_k8s_exists($client, 'apps/v1', 'Deployment', 'crossplane-system', 'crossplane')
```
