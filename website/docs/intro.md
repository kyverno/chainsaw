# Introduction

Chainsaw is a tool developed to **run end to end tests in [Kubernetes](https://kubernetes.io) clusters**.

It is meant to test [Kubernetes operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator) work as expected by running a sequence of test steps for:

- Creating resources
- Asserting operators react (or not) the way they should

## Why we made it ?

While developing [Kyverno](https://kyverno.io) we need to run end to end tests to make sure our admission controller works as expected.

!!! info "A typical Kyverno end to end test"

    [Kyverno](https://kyverno.io) can validate, mutate and generate resources based on policies installed in a cluster and a typical test is:

    1. Create a policy
    1. Create a resource
    1. Check that Kyverno acted as expected
    1. Cleanup and move to the next test

### From KUTTL to Chainsaw

We started with another tool called [KUTTL](https://kuttl.dev).

While KUTTL was a great tool to start with, we quickly identified some limitations and forked it to add the features we needed.

In the end we needed more flexibility than what KUTTL could offer and we started designing a new assertion model and at this point it was simpler to start a new tool from scratch than continuing making changes in our KUTTL fork.

The changes we were making was simply too large to have a chance to be incorporated upstream.

Chainsaw was born :tada:

### KUTTL compatibility

Fortunately, Chainsaw assertion model is compatible with KUTTL and therefore migration to Chainsaw is relatively easy.

We provide `chainsaw migrate kuttl config` and `chainsaw migrate kuttl tests` commands to assist with the migration, see [Migration from KUTTL](./more/kuttl-migration.md) for details.

## Use cases

Chainsaw is built with CI tools in mind - you only really need to download and execute it in your build script.

Installing it in your machine is entirely up to you, but still possible.

