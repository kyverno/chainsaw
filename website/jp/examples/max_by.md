```
max_by([{name: 'foo', count: `3`}, {name: 'bar', count: `7`}, {name: 'baz', count: `1`}], &count) == {name: 'bar', count: `7`}
```

```
max_by([{name: 'foo'}, {name: 'zoo'}, {name: 'bar'}], &name) == {name: 'zoo'}
```
