# time_utc

## Signature

`time_utc(string)`

## Description

Calcutes time in UTC from a given time in RFC 3339 format.

## Examples

```
time_utc('2024-01-15T10:30:00+02:00') == '2024-01-15T08:30:00Z'
```

```
time_utc('2024-01-15T10:30:00-05:00') == '2024-01-15T15:30:00Z'
```
