# regex_replace_all

## Signature

`regex_replace_all(string, string|number, string|number)`

## Description

Converts all parameters to string.

## Examples

```
regex_replace_all('[aeiou]', 'foobar', '*') == 'f**b*r'
```

```
regex_replace_all('[0-9]+', 'abc123def456', 'NUM') == 'abcNUMdefNUM'
```
