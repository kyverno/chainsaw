# time_add

## Signature

`time_add(string, string)`

## Description

Adds duration (second string) to a time value (first string).

## Examples

```
time_add('2024-01-01T00:00:00Z', '24h') == '2024-01-02T00:00:00Z'
```

```
time_add('2024-01-01T12:00:00Z', '30m') == '2024-01-01T12:30:00Z'
```
