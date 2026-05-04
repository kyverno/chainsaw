# trim_prefix

## Signature

`trim_prefix(string, string)`

## Description

Trims the second string prefix from the first string if the first string starts with the prefix.

## Examples

```
trim_prefix('foobar', 'foo') == 'bar'
```

```
trim_prefix('foobar', 'baz') == 'foobar'
```

```
trim_prefix('foobar', '') == 'foobar'
```
