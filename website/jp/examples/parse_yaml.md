```
parse_yaml('name: chainsaw') == {name: 'chainsaw'}
```

```
parse_yaml('enabled: true').enabled == `true`
```

```
parse_yaml('items:\n  - foo\n  - bar').items[0] == 'foo'
```
