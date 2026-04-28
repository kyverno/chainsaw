# min_by

## Signature

`min_by(array, expref)`

## Description

Returns the lowest found element using a custom expression to compute the associated value for each element in the input array.

## Examples

```
min_by([{name: 'foo', count: `3`}, {name: 'bar', count: `7`}, {name: 'baz', count: `1`}], &count) == {name: 'baz', count: `1`}
```

```
min_by([{name: 'foo'}, {name: 'zoo'}, {name: 'bar'}], &name) == {name: 'bar'}
```
