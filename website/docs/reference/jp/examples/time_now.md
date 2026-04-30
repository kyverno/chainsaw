# time_now

## Signature

`time_now()`

## Description

Returns current time in RFC 3339 format.

## Examples

!!! note "Non-deterministic output"

    The `time_now` function returns the current time and produces a different value each call.

```
# returns the current time in RFC 3339 format, e.g. '2024-01-15T10:30:00Z'
time_now()
```
