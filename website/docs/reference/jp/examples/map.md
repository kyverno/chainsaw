# map

## Signature

`map(expref, array)`

## Description

Transforms elements in a given array and returns the result.

## Examples

```
map(&to_upper(@), ['foo', 'bar']) == ['FOO', 'BAR']
```

```
map(&length(@), ['foo', 'foobar']) == [`3`, `6`]
```

```
map(&name, [{name: 'foo'}, {name: 'bar'}]) == ['foo', 'bar']
```
