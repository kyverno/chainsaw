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
1. [Apply operations](#apply)
1. [Assert operations](#assert)
1. [Error operations](#error)

### Delete operations   {#delete}

TODO

### Apply operations    {#apply}

TODO

### Assert operations   {#assert}

TODO

### Error operations    {#error}

TODO

## Cleanup

Unless configured differently, by default Chainsaw will automatically cleanup the resources it created after a test finishes.
Cleanup happens in reverse order of creation (created last, cleaned up first).

Note that Chainsaw performs a blocking deletion, that is, it will wait the resource is actually not present anymore in the cluster before proceeding with the next resource cleanup.

This is important, especially when the controller being tested makes use of `finalizers`.

## Namespaced resources

TODO
