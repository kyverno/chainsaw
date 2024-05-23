# Writing Chainsaw tests

This documentation focuses on providing a breakdown of the Chainsaw test structure and how to use it.

## What is a test?

To put it simply, a test can be represented as an ordered sequence of test steps.

In turn, a test step can be represented as an ordered sequence of operations.

To create a Chainsaw test all you need to do is to create one (or more) YAML file(s).

## Definition approach

Chainsaw supports two different test definition approaches:

!!! tip
    While Chainsaw supports two test definition approaches, we strongly recommend the [explicit one](./definition/explicit.md).

- [The explicit approach (strongly recommended)](./definition/explicit.md)
- [The conventional approach](./definition/conventional.md)

Each approach has pros and cons and choosing one over the other is not a trivial choice.

If you have no strong preference we encourage you to try the [explicit approach](./definition/explicit.md) first. The conventional one will only cover simple scenarios and will probably show its limitations as your tests become more complex.
