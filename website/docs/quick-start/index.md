# Chainsaw

Chainsaw is a tool primarily developed to **run end-to-end tests in [Kubernetes](https://kubernetes.io) clusters**.

It is meant to test [Kubernetes operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator) work as expected by running a sequence of steps and asserting various conditions.

## Why we made it?

While developing [Kyverno](https://kyverno.io) we need to run end-to-end tests to make sure our admission controller works as expected.

!!! info "A typical Kyverno end-to-end test"

    [Kyverno](https://kyverno.io) can validate, mutate and generate resources based on policies installed in a cluster and a typical test is:

    1. Create a policy
    1. Create a resource
    1. Check that Kyverno acted as expected
    1. Cleanup and move to the next test

## How to use it?

Chainsaw is built with CI tools in mind - you only really need to download and execute it in your build script.

However, [installing it on your local machine](./install.md) is entirely possible.
