# time_truncate

## Signature

`time_truncate(string, string)`

## Description

Returns the result of rounding time down to a multiple of duration.

## Examples

```
time_truncate('2024-01-15T10:37:45Z', '1h') == '2024-01-15T10:00:00Z'
```

```
time_truncate('2024-01-15T10:37:45Z', '24h') == '2024-01-15T00:00:00Z'
```
