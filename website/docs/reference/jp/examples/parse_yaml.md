# parse_yaml

## Signature

`parse_yaml(string)`

## Description

Decodes a valid YAML encoded string to the appropriate type provided it can be represented as JSON.

## Examples

```
parse_yaml('name: chainsaw') == {name: 'chainsaw'}
```

```
parse_yaml('enabled: true').enabled == `true`
```

```
parse_yaml('items:\n  - foo\n  - bar').items[0] == 'foo'
```
