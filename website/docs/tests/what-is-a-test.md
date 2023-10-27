# What is a test

In a nutshell, a test can be represented as **an ordered sequence of test steps**. Test steps within a test are run sequentially: **if any of the test steps fail, the entire test is considered failed**.

A test step can consist in **one or more operations**:

1. To delete resources present in a cluster
1. To create or update resources in a cluster
1. To assert one or more resources in a cluster meet the expectations (or the opposite)
1. To run arbitrary commands (will be supported soon)

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

1. [Delete operations](#delete)
2. [Apply operations](#apply)
3. [Assert operations](#assert)
4. [Error operations](#error)
5. [Exec operations](#exec)

### Delete operations   {#delete}

The delete operation allows you to specify resources that should be deleted from the Kubernetes cluster before a particular test step is executed.

#### Fields Description

The full structure of the `Delete` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test)

```yaml
Delete:
  - apiVersion: v1
    kind: Pod
    namespace: default
    name: my-test-pod
```

### Apply operations    {#apply}

The apply operation lets you define resources that should be applied to the Kubernetes cluster during the test step. These can be configurations, deployments, services, or any other Kubernetes resource.

The full structure of the `Apply` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test)

```yaml
Apply:
  - file: path/to/deployment.yaml
  - file: path/to/service.yaml
```

### Assert operations   {#assert}

The assert operation allows you to specify conditions that should hold true for a successful test. For example, after applying certain resources, you might want to ensure that a particular pod is running or a service is accessible.

The full structure of the `Assert` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test)

```yaml
Assert:
  - file: path/to/assertions.yaml  
```

### Error operations    {#error}

The error operation lets you define a set of expected errors for a test step. If any of these errors occur during the test, they are treated as expected outcomes. However, if an error that's not on this list occurs, it will be treated as a test failure.

The full structure of the `Error` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test)

```yaml
Error:
  - file: path/to/expected-errors.yaml
```

### Exec operations {#exec}

The `Exec` operation provides a means to either execute a specific command or run a script during the test step. For each `Exec` entry, you must specify either a `Command` or a `Script`, but not both.

The full structure of the `Exec` is documented [here](../apis/chainsaw.v1alpha1.md#chainsaw-kyverno-io-v1alpha1-Test)

#### Usage

To run a command:

```yaml
Exec:
  timeout: "5m"
  command:
    some-command-to-run-parameters
```

To execute a script:

```yaml
Exec:
  timeout: "10m"
  script:
    path/to/script.sh
```

> Make sure you're selecting either `Command` or `Script` for each `Exec` entry, and not both simultaneously.

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
