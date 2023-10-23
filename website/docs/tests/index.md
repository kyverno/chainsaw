# What is Chainsaw

## Overview

Chainsaw provides a **declarative approach** to testing production-grade Kubernetes [operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and controllers.

It provides a way to inject an operator (subject under test) during the setup and allows tests to be standard YAML files.

Test assertions are often partial YAML documents which assert the state defined is true.

## Motivation

Testing Kubernetes operators is not easy.

At [Kyverno](htpps://kyverno.io) we are building complex controllers and admission controllers that need to be deeply end to end tested.

It made sense to be able to declaratively create end to end tests for a couple of reasons:

- Easily translate a reported issue to an end to end test
- Add new end to end tests without writing a single line of code

The motivation is to leverage the existing Kubernetes eco-system for resource management (YAMLs) in a way to setup a test and as well as a way to assert state within the cluster.

## When would you use Chainsaw

The testing eco-system is vast and includes at a minimum low level unit tests, integration tests and end-to-end testing.

Chainsaw is built to support some kubernetes integration test scenarios and is most valuable as an end to end testing tool.

Chainsaw is great when you want to:

- Provide tests against your Custom Resource Definitions (CRDs)
- Inject a controller and assert states in a running cluster
- Run same tests against multiple implementations and multiple versions of Kubernetes clusters
