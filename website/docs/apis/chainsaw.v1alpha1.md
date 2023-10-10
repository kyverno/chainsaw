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
    



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>chainsaw.kyverno.io/v1alpha1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>Configuration</code></td></tr>
    
  
<tr><td><code>TypeMeta</code> <B>[Required]</B><br/>
<code>k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta</code>
</td>
<td>(Members of <code>TypeMeta</code> are embedded into this type.)
   <span class="text-muted">No description provided.</span></td>
</tr>
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<code>k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta</code>
</td>
<td>
   <span class="text-muted">No description provided.</span>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr><td><code>spec</code> <B>[Required]</B><br/>
<a href="#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec"><code>ConfigurationSpec</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span></td>
</tr>
</tbody>
</table>

## `ConfigurationSpec`     {#chainsaw-kyverno-io-v1alpha1-ConfigurationSpec}
    

**Appears in:**

- [Configuration](#chainsaw-kyverno-io-v1alpha1-Configuration)



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>duration</code> <B>[Required]</B><br/>
<code>int</code>
</td>
<td>
   <span class="text-muted">No description provided.</span></td>
</tr>
</tbody>
</table>
  