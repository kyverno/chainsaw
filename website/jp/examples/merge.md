```
merge({foo: 'bar'}, {baz: 'qux'}) == {foo: 'bar', baz: 'qux'}
```

```
merge({foo: 'bar'}, {foo: 'override'}) == {foo: 'override'}
```

```
merge({a: `1`}, {b: `2`}, {c: `3`}) == {a: `1`, b: `2`, c: `3`}
```
