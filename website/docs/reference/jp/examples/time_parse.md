# time_parse

## Signature

`time_parse(string, string)`

## Description

Changes a time value of a given layout to RFC 3339.

## Examples

```
time_parse('2006-01-02', '2024-03-15') == '2024-03-15T00:00:00Z'
```

```
time_parse('01/02/2006', '03/15/2024') == '2024-03-15T00:00:00Z'
```
