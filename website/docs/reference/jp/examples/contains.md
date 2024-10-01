# contains

## Signature

`contains(array|string, any)`

## Description

Returns `true` if the given subject contains the provided search value. If the subject is an array, this function returns `true` if one of the elements in the array is equal to the provided search value. If the provided subject is a string, this function returns `true` if the string contains the provided search argument.

## Examples

### With strings

```
contains('foobar', 'bar') == `true`
```

```
contains('foobar', 'not') == `false`
```

### With arrays

```
contains(['foo', 'bar'], 'bar') == `true`
```

```
contains(['foo', 'bar'], 'not') == `true`
```
