# Test spec


## Supported elements

| Field | Default | Description |
|---|---|---|
| `namespace` | | Namespace determines whether the test should run in a random ephemeral namespace or not. |
| `namespaceTemplate` | | NamespaceTemplate defines a template to create the test namespace. |
| `timeouts` | | Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis. |
| `steps` | | Steps defining the test. |
| `clusters` | | Clusters holds a registry to clusters to support multi-cluster tests. |
| `cluster` | | Cluster defines the target cluster (default cluster will be used if not specified and/or overridden). |
| `bindings` | | Bindings defines additional binding key/values. |
| `catch` | | Catch defines what the steps will execute when an error happens. This will be combined with catch handlers defined at the step level. |
| `template` | | Template determines whether resources should be considered for templating. |
| `concurrent` | | Concurrent determines whether the test should run concurrently with other tests. |
| `skip` | | Skip determines whether the test should skipped. |
| `skipDelete` | | SkipDelete determines whether the resources created by the test should be deleted after the test is executed. |
| `forceTerminationGracePeriod` | | ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments. |
| `delayBeforeCleanup` | | DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts. |
| `deletionPropagationPolicy` | | DeletionPropagationPolicy decides if a deletion will propagate to the dependents of the object, and how the garbage collector will handle the propagation. Overrides the deletion propagation policy set in the Configuration. |
| `description` | | Description contains a description of the test. |

### Namespace

### Namespace template

### Timeouts

### Steps

### Clusters

### Cluster

### Bindings

### Catch

### Template

### Concurrency

