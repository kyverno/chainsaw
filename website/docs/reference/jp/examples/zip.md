# zip

## Signature

`zip(array, array)`

## Description

Accepts one or more arrays as arguments and returns an array of arrays in which the i-th array contains the i-th element from each of the argument arrays. The returned array is truncated to the length of the shortest argument array.

## Examples

```
zip(['a', 'b'], [`1`, `2`]) == [['a', `1`], ['b', `2`]]
```
