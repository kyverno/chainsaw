# random

## Signature

`random(string)`

## Description

Generates a random sequence of characters.

## Examples

!!! note "Non-deterministic output"

    The `random` function generates a different value each call. The examples below illustrate the pattern, not an exact return value.

```
# generates a random alphanumeric string
random('[a-z0-9]{8}')
```
