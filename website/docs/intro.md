# Introduction

Chainsaw is a tool developed to run end to end tests in [Kubernetes](https://kubernetes.io) clusters.

It is meant to test [Kubernetes operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator) work as expected by running a sequence of test steps, creating resources and asserting operators react (or not) in the expected way.

We use it internally to test [Kyverno](https://kyverno.io).

## Why we made it ?

While developing [Kyverno](https://kyverno.io) we needed to run end to end tests to make sure our admission controller worked as expected.

[Kyverno](https://kyverno.io) can validate, mutate and generate resources based on policies installed in a cluster and a typical test is:
1. Create a policy
1. Create a resource
1. Check that Kyverno acted as expected
1. Cleanup and move to the next test

We started with another tool called [Kuttl](https://kuttl.dev). While [Kuttl](https://kuttl.dev) is great we identified some limitations and forked the tool to add the features we needed.

At some point we needed more flexibility than what [Kuttl](https://kuttl.dev) offered and we designed a new assertion model.

This was simpler to start from scratch than continuing making changes in our [Kuttl](https://kuttl.dev) fork.

Migrating from [Kuttl](https://kuttl.dev) to Chainsaw should be easy, the assertion model used in Kuttl is compatible with Chainsaw.

## Use cases

Chainsaw is built with CI tools in mind - you only really need to download and execute it in your build script.

Installing it in your machine is entirely up to you, but still possible.

