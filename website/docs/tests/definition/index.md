# Definition approach

## Writing tests

Chainsaw supports two different test definition approaches:

!!! tip
    While Chainsaw supports two test definition approaches, we strongly recommend the [explicit one](./explicit.md).

- [The explicit approach (strongly recommended)](./explicit.md)
- [The conventional approach](./conventional.md)

Each approach has pros and cons and choosing one over the other is not a trivial choice.

If you have no strong preference we encourage you to try the [explicit approach](./explicit.md) first. The conventional one will only cover simple scenarios and will probably show its limitations as your tests become more complex.

---

## KUTTL

The commands below can be used to help migrating from KUTTL to Chainsaw:

- [`chainsaw migrate kuttl config`](../../reference/commands/chainsaw_migrate_kuttl_config.md)

    Migrate a KUTTL test suite to the corresponding Chainsaw configuration

- [`chainsaw migrate kuttl tests`](../../reference/commands/chainsaw_migrate_kuttl_tests.md)

    Migrate KUTTL tests to the corresponding Chainsaw tests
