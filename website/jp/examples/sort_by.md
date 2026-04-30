```
sort_by([{name: 'foo', count: `3`}, {name: 'bar', count: `1`}, {name: 'baz', count: `2`}], &count) == [{name: 'bar', count: `1`}, {name: 'baz', count: `2`}, {name: 'foo', count: `3`}]
```

```
sort_by([{name: 'foo'}, {name: 'bar'}, {name: 'baz'}], &name) == [{name: 'bar'}, {name: 'baz'}, {name: 'foo'}]
```
