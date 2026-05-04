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
