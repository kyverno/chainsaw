# semver_compare

## Signature

`semver_compare(string, string)`

## Description

Compares two strings which comply with the semantic versioning schema and outputs a boolean response as to the position of the second relative to the first.

## Examples

```
semver_compare('1.0.0', '>1.0.0') == `false`
```

```
semver_compare('1.2.0', '>1.0.0') == `true`
```

```
semver_compare('1.0.0', '>=1.0.0') == `true`
```

```
semver_compare('2.0.0', '<1.5.0') == `false`
```

```
semver_compare('1.0.0', '=1.0.0') == `true`
```
