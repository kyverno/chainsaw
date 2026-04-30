# time_before

## Signature

`time_before(string, string)`

## Description

Checks if a time is before another time, both in RFC3339 format.

## Examples

```
time_before('2024-01-01T00:00:00Z', '2024-01-02T00:00:00Z') == `true`
```

```
time_before('2024-01-02T00:00:00Z', '2024-01-01T00:00:00Z') == `false`
```

```
time_before('2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z') == `false`
```
