# What is a test

To put it simply, a test can be represented as **an ordered sequence of [test steps](../steps/index.md)**.

Test steps within a test are run sequentially: **if any of the test steps fail, the entire test is considered failed**.

A test step can consist in **one or more operations**:

- To **delete** resources present in a cluster
- To **create** or update resources in a cluster
- To **assert** one or more resources in a cluster meet the expectations (or the opposite)
- To run arbitrary **commands** or **scripts**

## Different syntaxes

Chainsaw supports two different test definition syntaxes. Each syntax has pros and cons, see the descriptions below for more details on advantages and pitfalls.

### [Manifests based syntax](./manifests-based.md)

This is the simplest and less verbose syntax.

You provide bare Kubernetes resource manifests and Chainsaw will use those manifests to create, update, or assert expectations against a cluster.

While this syntax is extremly simple, not all operations are supported. For example, `delete`, `command`, `script` and `sleep` operations are not.

Another strong limitation is that it is not possible to specify additional configuration per test, step or operation.

Finally, this syntax relies heavily on file naming conventions, it can be error prone and makes it hard to reuse files across multiple tests.

### [`Test` based syntax](./test-based.md)

The `Test` based syntax is a more verbose and explicit syntax.

It does not rely on file naming conventions to determine test steps order and allows to easily reuse files accross multiple tests.

This syntax also comes with no limitations to provide additional configuration at the test, step or operation level.

!!! tip "Making a choice"

    Choosing one syntax over the other is not a trivial choice, every one will have its own preference and/or constraints.

    It's usually easier to start with the manifests based syntax.
    However, as test suites grow and tests become more complex, it is often necessary to configure options on a per test, step or operation basis and the `Test` based syntax becomes necessary.

    Fortunately Chainsaw has a command to [automatically migrate](../commands/chainsaw_migrate_tests.md) from manifest based to `Test` based syntax.

## Namespaced resources

Kubernetes organizes resources into two primary scopes: namespaced and cluster-scoped.

While namespaced resources belong to a specific namespace, cluster-scoped resources span across the entire Kubernetes cluster.

### Purpose of the Namespacer Interface

The [Namespacer interface](https://github.com/kyverno/chainsaw/blob/main/pkg/runner/namespacer/namespacer.go) ensures automated and consistent namespace assignment to Kubernetes resources.

- **Automated Namespacing**

    Automatically assign namespaces to resources that don't have one.

- **Ephemeral Namespaces**

    Handles temporary namespaces for specific tasks.

## Cleanup

Unless configured differently, by default Chainsaw will automatically cleanup the resources it created after a test finishes.
Cleanup happens in reverse order of creation (created last, cleaned up first).

Note that Chainsaw performs a blocking deletion, that is, it will wait the resource is actually not present anymore in the cluster before proceeding with the next resource cleanup.

This is important, especially when the controller being tested makes use of `finalizers`.

!!! tip "Overriding cleanup timeout"

    A global cleanup timeout can be defined at the configuration level or using command line flags.

    It can also be overriden on a per test or per test step basis but not at the operation level.

## Parallel Execution of Tests

While Chainsaw ensures that the steps within a test are executed sequentially, it is designed to run multiple tests in parallel to each other. This parallel execution helps in significantly reducing the overall time required to run an entire suite of tests, making the testing process more efficient, especially in scenarios with a large number of tests.

!!! tip "Parallel execution of tests"

    By default, Chainsaw will run tests in parallel.

    This can be configured at the configuration level or using command line flags.
