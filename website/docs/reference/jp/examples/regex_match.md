# regex_match

## Signature

`regex_match(string, string|number)`

## Description

First string is the regular exression which is compared with second input which can be a number or string.

## Examples

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
