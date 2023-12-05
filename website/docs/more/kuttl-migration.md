# Migration from KUTTL

## Overview

The `chainsaw migrate kuttl tests` and `chainsaw migrate kuttl config` commands are designed for the migration of KUTTL tests to Chainsaw.

- [`chainsaw migrate kuttl config`](#migrate-configuration)

    migrates a KUTTL `TestSuite` to the corresponding Chainsaw `Configuration`

- [`chainsaw migrate kuttl tests`](#migrate-tests)

    migrates KUTTL tests to the corresponding Chainsaw `Test`s

!!! tip "Reference documentation"

    You can view the full command documentation [here](../commands/chainsaw_migrate_kuttl.md).

## Examples

### Migrate tests

The command below will migrate KUTTL tests to Chainsaw and overwrite original files with converted ones.

```bash
chainsaw migrate kuttl tests path/to/kuttl/tests --save --cleanup
```

This will generate a `chainsaw-test.yaml` for every KUTTL test discovered.

### Migrate configuration

The command below will migrate a KUTTL test suite file to the corresponding Chainsaw `Configuration`.

```bash
chainsaw migrate kuttl config path/to/kuttl/testsuite --save --cleanup
```

This will generate a `.chainsaw.yaml` configuration file.

