# Chainsaw

[![Lint](https://github.com/kyverno/chainsaw/actions/workflows/lint.yaml/badge.svg)](https://github.com/kyverno/chainsaw/actions/workflows/lint.yaml)
[![Tests](https://github.com/kyverno/chainsaw/actions/workflows/tests.yaml/badge.svg)](https://github.com/kyverno/chainsaw/actions/workflows/tests.yaml)
[![Code QL](https://github.com/kyverno/chainsaw/actions/workflows/codeql.yaml/badge.svg)](https://github.com/kyverno/chainsaw/actions/workflows/codeql.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyverno/chainsaw)](https://goreportcard.com/report/github.com/kyverno/chainsaw)
[![License: Apache-2.0](https://img.shields.io/github/license/kyverno/chainsaw?color=blue)](https://github.com/kyverno/chainsaw/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kyverno/chainsaw/branch/main/graph/badge.svg)](https://app.codecov.io/gh/kyverno/chainsaw/branch/main)

Chainsaw provides a declarative approach to test [Kubernetes](https://kubernetes.io) operators and controllers.

While Chainsaw is designed for testing operators and controllers, it can declaratively test any kubernetes objects.

Chainsaw is an open source tool that was initially developped for defining and running [Kyverno](https://kyverno.io) end to end tests.
