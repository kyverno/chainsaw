---
title: chainsaw (v1alpha1)
content_type: tool-reference
package: chainsaw.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the v1alpha1 API group</p>


## Resource Types 


- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)
  
## `Configuration`     {#chainsaw-kyverno-io-v1alpha1-Configuration}

<p>Configuration is the resource that contains the configuration used to run tests.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | `chainsaw.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | `Configuration` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  | <p>Standard object's metadata.</p> |
| `spec` | [`ConfigurationSpec`](#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec) | :white_check_mark: | <p>Configuration spec.</p> |

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}

**Appears in:**
    
- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)

<p>ConfigurationSpec contains the configuration used to run tests.</p>


| Field | Type | Required | Description |
|---|---|---|---|
| `timeout` | `int` | :white_check_mark: | *No description provided.* |

  