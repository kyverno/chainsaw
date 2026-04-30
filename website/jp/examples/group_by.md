```
group_by([{name: 'foo', type: 'a'}, {name: 'bar', type: 'b'}, {name: 'baz', type: 'a'}], &type) == {a: [{name: 'foo', type: 'a'}, {name: 'baz', type: 'a'}], b: [{name: 'bar', type: 'b'}]}
```
