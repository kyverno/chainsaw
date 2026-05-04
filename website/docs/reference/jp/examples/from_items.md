# from_items

## Signature

`from_items(array[array])`

## Description

Returns an object from the provided array of key value pairs. This function is the inversed of the `items()` function.

## Examples

```
from_items([['foo', 'bar'], ['baz', 'qux']]) == {foo: 'bar', baz: 'qux'}
```

```
from_items([['a', `1`], ['b', `2`]]) == {a: `1`, b: `2`}
```
