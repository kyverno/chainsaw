# Multi-cluster options

Multi-cluster options contain the configuration of additional clusters.

## Supported elements

Every cluster is registered by name and supports the following elements:

| Element | Default | Description |
|---|---|---|
| `kubeconfig` | `string` | Kubeconfig is the path to the referenced file. |
| `context` | `string` | Context is the name of the context to use. |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: custom-config
spec:
  clusters:
    # this cluster will use the default (current) context
    # configured in the kubeconfig file
    cluster-1:
      kubeconfig: /path/to/kubeconfig-1
    # this cluster will use the context named `context-2`
    # in the kubeconfig file
    cluster-2:
      kubeconfig: /path/to/kubeconfig-2
      context: context-2
```

### With flags

!!! note
    The `--cluster` flag can appear multiple times and is expected to come in the following format:
    
    `--cluster cluster-name=/path/to/kubeconfig[:context-name]`.

```bash
chainsaw test                                               \
    --cluster cluster-1=/path/to/kubeconfig-1               \
    --cluster cluster-2=/path/to/kubeconfig-2:context-2
```
