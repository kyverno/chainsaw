# Writing Chainsaw tests

This documentation focuses on providing a breakdown of the Chainsaw test structure and how to use it.

## What is a test?

To put it simply, a test can be represented as an ordered sequence of test steps.

In turn, a test step can be represented as an ordered sequence of operations.

- When an operation fails the test is considered failed
- If all operations succeed the test is considered successful

## Definition approach

Chainsaw supports two different test definition approaches:

!!! tip
    While Chainsaw supports two test definition approaches, we strongly recommend the [explicit one](./explicit.md).

- [The explicit approach (strongly recommended)](./explicit.md)
- [The conventional approach](./conventional.md)

## General concepts

The concepts below are at the heart of Chainsaw:

- [Inheritance](../general/inheritance.md)
- [Test namespace](../general/namespace.md)
- [Bindings](../general/bindings.md)
- [Templating](../general/templating.md)
- [Outputs](../general/outputs.md)
- [References](../general/references.md)
- [Operation checks](../general/checks.md)

## Test and Step specs

Browse the test and step specs to learn all the details and options:

- [Test spec](./spec/index.md)
- [Test step spec](../step/index.md)
