# Collectors

## Purpose

The purpose of collectors is to collect certain information about the outcome of a step should it fail (in the case of `catch`) or at the end of the step (in the case of `finally`).

The ultimate goal of collectors is to gather information about the failure of a step and therefore help understand what caused it to fail.

A test step can have an arbitrary number of collectors.

## Available collectors

- [Pod logs](./pod-logs.md)
- [Events](./events.md)
- [Get](./get.md)
- [Describe](./describe.md)
