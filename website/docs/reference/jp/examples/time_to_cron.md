# time_to_cron

## Signature

`time_to_cron(string)`

## Description

Converts a time (RFC 3339) to a cron expression (string).

## Examples

```
time_to_cron('2024-01-15T10:30:00Z') == '30 10 15 1 1'
```
