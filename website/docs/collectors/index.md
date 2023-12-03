# Collectors

The purpose of collectors is to collect certain information about the outcome of a step should it fail (in the case of `catch`) or at the end of the step (in the case of `finally`).

The ultimate goal of collectors is to gather information about the failure of a step and therefore help understand what caused it to fail.

A test step can have an arbitrary number of collectors.

## Collectors lifecycle

Collectors are executed **after** the step operations, and **before** the step cleanup happens.

!!! info "Catch lifecycle"

    1. The step starts executing
    1. An operation fails or all operations in the step terminate
    1. Catch operations and collectors are executed (**if the step failed**)
    1. Finally operations and collectors are executed (**in all cases**)
    1. The step cleanup executes

This is important that collectors run before cleanup so that they have a chance to collect logs from pods responsible for the step failure.

## Available collectors

- [Pod logs](./pod-logs.md)
- [Events](./events.md)
