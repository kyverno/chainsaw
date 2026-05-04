# items

## Signature

`items(object|array, string, string)`

## Description

Converts a map or array to an array of objects where each key:value is an item in the array.

## Examples

### With an object

```
items({foo: 'bar', baz: 'qux'}, 'key', 'value') == [{key: 'baz', value: 'qux'}, {key: 'foo', value: 'bar'}]
```

### With an array

```
items(['foo', 'bar'], 'key', 'value') == [{key: `0`, value: 'foo'}, {key: `1`, value: 'bar'}]
```
