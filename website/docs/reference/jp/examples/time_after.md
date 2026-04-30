# time_after

## Signature

`time_after(string, string)`

## Description

Checks if a time is after another time, both in RFC3339 format.

## Examples

```
time_after('2024-01-02T00:00:00Z', '2024-01-01T00:00:00Z') == `true`
```

```
time_after('2024-01-01T00:00:00Z', '2024-01-02T00:00:00Z') == `false`
```

```
time_after('2024-01-01T00:00:00Z', '2024-01-01T00:00:00Z') == `false`
```
