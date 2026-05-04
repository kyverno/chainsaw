# max

## Signature

`max(array[number]|array[string])`

## Description

Returns the highest found element in the provided array argument. An empty array will produce a return value of null.

## Examples

### With numbers

```
max([`1`, `5`, `3`]) == `5`
```

### With strings

```
max(['b', 'a', 'c']) == 'c'
```

### With an empty array

```
max(`[]`) == null
```
