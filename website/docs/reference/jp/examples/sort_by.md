# sort_by

## Signature

`sort_by(array, expref)`

## Description

This function accepts an array argument and returns the sorted elements as an array using a custom expression to compute the associated value for each element.

## Examples

```
sort_by([{name: 'foo', count: `3`}, {name: 'bar', count: `1`}, {name: 'baz', count: `2`}], &count) == [{name: 'bar', count: `1`}, {name: 'baz', count: `2`}, {name: 'foo', count: `3`}]
```

```
sort_by([{name: 'foo'}, {name: 'bar'}, {name: 'baz'}], &name) == [{name: 'bar'}, {name: 'baz'}, {name: 'foo'}]
```
