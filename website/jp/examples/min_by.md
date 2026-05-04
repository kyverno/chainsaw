```
min_by([{name: 'foo', count: `3`}, {name: 'bar', count: `7`}, {name: 'baz', count: `1`}], &count) == {name: 'baz', count: `1`}
```

```
min_by([{name: 'foo'}, {name: 'zoo'}, {name: 'bar'}], &name) == {name: 'bar'}
```
