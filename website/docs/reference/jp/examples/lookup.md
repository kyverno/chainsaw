# lookup

## Signature

`lookup(object|array, string|number)`

## Description

Returns the value corresponding to the given key/index in the given object/array.

## Examples

### With an object

```
lookup({foo: 'bar', baz: 'qux'}, 'foo') == 'bar'
```

```
lookup({foo: 'bar', baz: 'qux'}, 'missing') == null
```

### With an array

```
lookup(['foo', 'bar', 'baz'], `1`) == 'bar'
```
