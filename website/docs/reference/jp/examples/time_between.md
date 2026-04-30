# time_between

## Signature

`time_between(string, string, string)`

## Description

Checks if a time is between a start and end time, all in RFC3339 format.

## Examples

```
time_between('2024-06-15T00:00:00Z', '2024-01-01T00:00:00Z', '2024-12-31T00:00:00Z') == `true`
```

```
time_between('2023-06-15T00:00:00Z', '2024-01-01T00:00:00Z', '2024-12-31T00:00:00Z') == `false`
```

```
time_between('2025-01-01T00:00:00Z', '2024-01-01T00:00:00Z', '2024-12-31T00:00:00Z') == `false`
```
