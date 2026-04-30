# regex_replace_all_literal

## Signature

`regex_replace_all_literal(string, string|number, string|number)`

## Description

Converts all parameters to string.

## Examples

```
regex_replace_all_literal('[aeiou]', 'foobar', '$0$0') == 'f$0$0$0$0b$0$0r'
```

```
regex_replace_all_literal('[0-9]+', 'abc123def456', '$NUM') == 'abc$NUMdef$NUM'
```
