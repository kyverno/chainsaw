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

## Use cases

Chainsaw is built with CI tools in mind - you only really need to download and execute it in your build script.

Installing it in your machine is entirely up to you, but still possible.

