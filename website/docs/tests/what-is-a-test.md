# What is a test

In a nutshell, a test can be represented as **an ordered sequence of test steps**. Test steps within a test are run sequentially: **if any of the test steps fail, the entire test is considered failed**.

A test step can consist in **one or more operations**:

- To delete resources present in a cluster
- To create or update resources in a cluster
- To assert one or more resources in a cluster meet the expectations (or the opposite)
- To run arbitrary commands (will be supported soon)

## Different syntaxes

Chainsaw supports multiple test definition mechanisms. Under certain circumstances, some operations might not be directly available and it might be possible to combine multiple syntaxes together to assemble a test.

1. **Manifests based syntax**:
This is the simplest and less verbose supported syntax, you provide bare Kubernetes resource manifests and Chainsaw will use those manifests to create, update, or assert expectations against a cluster.
While this syntax is simple, it doesn't support deletion for example and doesn't allow specifying additional configuration per test or step.
It also relies a lot on file naming conventions and makes it hard to reuse files across multiple tests.

1. **`TestStep`s based syntax**:
This syntax is more verbose than the first one but offers the flexibility to provide additional configuration per test step.
It also makes it easier to reuse files accross multiple tests but still relies on the same file naming conventions as the manifests based syntax.
On the other hand it doesn't suffer the unsupported deletion limitation and can be combined with manifests based syntax when defining a step.

1. **`Test` based syntax**:
The `Test` based syntax is the more verbose and explicit syntax. It does not rely on file naming conventions, it makes it easy to reuse files accross multiple tests, and offers the flexibility to provide additional configuration at both the test level and test step level.
It supports all kind of operations.

Choosing one syntax over the other is not a trivial choice, every one will have its own preference and/or constraints.
If you feel more comfortable with explicit declarative models, the recommandation would be to use `Test` based syntax.
If you don't mind relying on file naming conventions and don't need to reuse files across multiple tests manifests based syntax or `TestStep`s based syntax is a good choice.
However, using `TestStep`s based syntax only is debatable. In this case, `Test` based syntax could be a simpler choice with more flexibility.

## Operations

Chainsaw supports the following operations, executed in this specific order for a given test step:

- [Delete operations](#delete)
- [Apply operations](#apply)
- [Create operations](#create)
- [Assert operations](#assert)
- [Error operations](#error)
- [Command operations](#command)
- [Script operations](#script)

### Common fields

All operations share some configuration fields.

- **Timeout:** A timeout for the operation.
- **ContinueOnError:** Determines whether a test step should continue or not in case the operation was not successful.
  Even if the test continues executing, it will still be reported as failed.

The full structure of the `Operation` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Operation).

### Delete operations   {#delete}

The delete operation allows you to specify resources that should be deleted from the Kubernetes cluster before a particular test step is executed.

The full structure of the `Delete` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Delete).

```yaml
delete:
  apiVersion: v1
  kind: Pod
  namespace: default
  name: my-test-pod
```

### Apply operations    {#apply}

The apply operation lets you define resources that should be applied to the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

The full structure of the `Apply` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Apply).

```yaml
apply:
  file: path/to/deployment.yaml
```

### Create operations    {#creatye}

The create operation lets you define resources that should be created in the Kubernetes cluster during the test step.
These can be configurations, deployments, services, or any other Kubernetes resource.

If the resource to be created already exists in the cluster, the step will fail.

The full structure of the `Create` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Create).

```yaml
create:
  file: path/to/deployment.yaml
```

### Assert operations   {#assert}

The assert operation allows you to specify conditions that should hold true for a successful test. For example, after applying certain resources, you might want to ensure that a particular pod is running or a service is accessible.

The full structure of the `Assert` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Assert).

```yaml
assert:
  file: path/to/assertions.yaml  
```

### Error operations    {#error}

The error operation lets you define a set of expected errors for a test step. If any of these errors occur during the test, they are treated as expected outcomes. However, if an error that's not on this list occurs, it will be treated as a test failure.

The full structure of the `Error` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Error).

```yaml
error:
  file: path/to/expected-errors.yaml
```

### Command operations {#command}

The `Command` operation provides a means to execute a specific command during the test step.

The full structure of the `Command` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Command).

```yaml
command:
  entrypoint: echo
  args:
  - hello chainsaw
```

### Script operations {#script}

The `Script` operation provides a means to run a script during the test step.

The full structure of the `Script` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Script).

```yaml
script:
  content: |
    echo "hello chainsaw"
```

## Cleanup

Unless configured differently, by default Chainsaw will automatically cleanup the resources it created after a test finishes.
Cleanup happens in reverse order of creation (created last, cleaned up first).

Note that Chainsaw performs a blocking deletion, that is, it will wait the resource is actually not present anymore in the cluster before proceeding with the next resource cleanup.

This is important, especially when the controller being tested makes use of `finalizers`.

## Namespaced resources

Kubernetes organizes resources into two primary scopes: namespaced and cluster-scoped. While namespaced resources belong to a specific namespace, cluster-scoped resources span across the entire Kubernetes cluster.

### Purpose of the Namespacer Interface

The [Namespacer interface](https://github.com/kyverno/chainsaw/blob/main/pkg/runner/namespacer/namespacer.go#L8) ensures automated and consistent namespace assignment to Kubernetes resources.

- **Automated Namespacing**: Automatically assign namespaces to resources that don't have one.
- **Ephemeral Namespaces**: Handles temporary namespaces for specific tasks.
