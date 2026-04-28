### With an object

```
items({foo: 'bar', baz: 'qux'}, 'key', 'value') == [{key: 'baz', value: 'qux'}, {key: 'foo', value: 'bar'}]
```

### With an array

```
items(['foo', 'bar'], 'key', 'value') == [{key: `0`, value: 'foo'}, {key: `1`, value: 'bar'}]
```
