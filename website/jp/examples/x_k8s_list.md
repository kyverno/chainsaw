!!! info "Clustered resources"

    For clustered resources, you can leave the namespace empty `''`.

```
# `$client` is a binding pointing to a Kubernetes client

x_k8s_list($client, 'apps/v1', 'Deployment', 'crossplane-system')
```