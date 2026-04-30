# object_from_lists

## Signature

`object_from_lists(array, array)`

## Description

Converts a pair of lists containing keys and values to an object.

## Examples

```
object_from_lists(['foo', 'bar', 'baz'], [`1`, `2`, `3`]) == {foo: `1`, bar: `2`, baz: `3`}
```

```
object_from_lists(['a', 'b'], ['x', 'y']) == {a: 'x', b: 'y'}
```
