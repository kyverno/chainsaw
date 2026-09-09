# Ecosystem & Related Projects

Chainsaw is designed to be modular, extensible, and reusable. Several of Chainsaw's core packages—such as resource templating, JMESPath expressions, and resource matching/validation—are exposed as public Go modules that can be integrated into other projects.

This page lists community-driven libraries, tools, and integrations that build on top of or complement Chainsaw.

---

## Sawchain

[Sawchain](https://github.com/guidewire-oss/sawchain) is an open-source Go testing library developed by Guidewire. It bridges the gap between Chainsaw's declarative testing power and programmatic Go test suites.

### How it works with Chainsaw
Instead of running as a standalone command-line runner, Sawchain integrates Chainsaw's core packages directly into Go's standard `testing.TB` framework (including Ginkgo and Gomega). It leverages Chainsaw's:
- **Resource templating** for dynamic test inputs.
- **JMESPath expressions** for rich queries.
- **Resource matching and assertion engine** for declarative checking.

### Key Features
- **Gomega Matchers**: Provides Gomega matchers for declarative Kubernetes resource verification inside Go test assertions.
- **Controller-runtime Integration**: Integrates directly with controller-runtime clients for smooth interaction with Kubernetes APIs.
- **Ergonomic DSL**: Allows developer-friendly, fluent test configurations directly in Go codebase.

### Resources
- [Sawchain GitHub Repository](https://github.com/guidewire-oss/sawchain)
- [Kubernetes Deserves Better Tests: Meet Sawchain (Deep Dive Blog)](https://medium.com/guidewire-engineering-blog/kubernetes-deserves-better-tests-meet-sawchain-2e8c8f750501)

---

!!! tip "Have a project built on Chainsaw?"

    We welcome community contributions! Please open a pull request to add your related project or integration to this list.
