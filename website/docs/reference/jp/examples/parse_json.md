# parse_json

## Signature

`parse_json(string)`

## Description

Decodes a valid JSON encoded string to the appropriate type. Opposite of `to_string` function.

## Examples

```
parse_json('{"name":"chainsaw","version":"1.0.0"}') == `{"name":"chainsaw","version":"1.0.0"}`
```

```
parse_json('{"enabled":true}').enabled == `true`
```

```
parse_json('[1,2,3]')[1] == `2`
```

```
parse_json('{"items":[{"id":1},{"id":2}]}').items[0].id == `1`
```
