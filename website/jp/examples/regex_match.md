```
regex_match('^foo.*', 'foobar') == `true`
```

```
regex_match('^foo.*', 'barfoo') == `false`
```

```
regex_match('^[0-9]+$', '12345') == `true`
```

```
regex_match('^[0-9]+$', '123abc') == `false`
```
