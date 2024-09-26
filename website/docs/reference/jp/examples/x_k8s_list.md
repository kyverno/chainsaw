# x_k8s_list

## Signature

`x_k8s_list(any, string, string, string)`

## Description

Lists resources from a Kubernetes cluster.

## Examples

!!! info "Clustered resources"

    For clustered resources, you can leave the namespace empty `''`.

```
# `$client` is a binding pointing to a Kubernetes client

x_k8s_list($client, 'apps/v1', 'Deployment', 'crossplane-system')
```
