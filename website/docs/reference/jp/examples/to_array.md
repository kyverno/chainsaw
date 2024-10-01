# to_array

## Signature

`to_array(any)`

## Description

Returns a one element array containing the passed in argument, or the passed in value if it's an array.

## Examples

```
to_array(`true`) == [`true`]
```

```
to_array([`10`, `15`, `20`]) == [`10`, `15`, `20`]
```

```
to_array(`[]`) == []
```
