# equal_fold

## Signature

`equal_fold(string, string)`

## Description

Allows comparing two strings for equivalency where the only differences are letter cases.

## Examples

```
equal_fold('Hello', 'hello') == `true`
```

```
equal_fold('FOOBAR', 'foobar') == `true`
```

```
equal_fold('foo', 'bar') == `false`
```
