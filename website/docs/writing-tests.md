# What is Chainsaw

## Overview

Chainsaw provides a **declarative approach** to testing production-grade Kubernetes [operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and controllers.

It provides a way to inject an operator (*subject under test*) during the setup and allows tests to be standard YAML files.

!!! info ""
    In Chainsaw **everything is YAML**. We use YAML to manipulate resources, define cluster state assertions, describe commands to run, etc...
    Everything is **declarative**.

## Motivation

Testing Kubernetes operators is not easy.

At [Kyverno](htpps://kyverno.io) we are building complex controllers and admission controllers that need to be deeply end to end tested.

!!! tip "Advantages of a declarative approach"

    It made sense to be able to declaratively create end to end tests for a couple of reasons:

    - Easily translate a reported issue to an end to end test
    - Add new end to end tests without writing a single line of code

## When would you use Chainsaw

The testing eco-system is vast and includes at a minimum low level unit tests, integration tests and end-to-end testing.

Chainsaw is built to support some kubernetes integration test scenarios and is most valuable as an **end to end testing tool**.

!!! tip "Chainsaw is great for"
    - Provide tests against your Custom Resource Definitions (CRDs)
    - Inject a controller and assert states in a running cluster
    - Run same tests against multiple implementations and multiple versions of Kubernetes clusters
