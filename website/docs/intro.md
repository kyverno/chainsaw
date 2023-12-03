# Introduction

Chainsaw is a tool developed to **run end to end tests in [Kubernetes](https://kubernetes.io) clusters**.

!!! info "How it works"

    It is meant to test [Kubernetes operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator) work as expected by running a sequence of test steps for:

    - creating resources
    - and asserting operators react (or not) in the expected way.

## Why we made it ?

While developing [Kyverno](https://kyverno.io) we needed to run end to end tests to make sure our admission controller worked as expected.

!!! info "A typical Kyverno end to end test"

    [Kyverno](https://kyverno.io) can validate, mutate and generate resources based on policies installed in a cluster and a typical test is:

    1. Create a policy
    1. Create a resource
    1. Check that Kyverno acted as expected
    1. Cleanup and move to the next test

We started with another tool called [KUTTL](https://kuttl.dev). While [KUTTL](https://kuttl.dev) is great we identified some limitations and forked the tool to add the features we needed.

But we ended up needing more flexibility than what [KUTTL](https://kuttl.dev) could offer and we started designing a new assertion model.

At this point it was simpler to start a new tool from scratch than continuing making changes in our [KUTTL](https://kuttl.dev) fork.

!!! tip "KUTTL compatibility"
    Migrating from [KUTTL](https://kuttl.dev) to Chainsaw should be easy, the assertion model used in KUTTL is compatible with Chainsaw. We provide `chainsaw kuttl migrate config` and `chainsaw kuttl migrate tests` commands to help with the migration, see [Migration from KUTTL](./more/kuttl-migration.md) for details.

## Use cases

Chainsaw is built with CI tools in mind - you only really need to download and execute it in your build script.

Installing it in your machine is entirely up to you, but still possible.

