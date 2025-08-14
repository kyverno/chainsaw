# wildcard

## Signature

`wildcard(string, string)`

## Description

Compares a wildcard pattern with a given string and returns if they match or not.

## Examples

```
wildcard('foo*', 'foobar') == `true`
```

```
wildcard('fooba?', 'foobar') == `true`
```
