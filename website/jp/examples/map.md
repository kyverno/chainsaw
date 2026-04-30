```
map(&to_upper(@), ['foo', 'bar']) == ['FOO', 'BAR']
```

```
map(&length(@), ['foo', 'foobar']) == [`3`, `6`]
```

```
map(&name, [{name: 'foo'}, {name: 'bar'}]) == ['foo', 'bar']
```
