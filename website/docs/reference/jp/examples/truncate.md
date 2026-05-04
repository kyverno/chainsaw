# truncate

## Signature

`truncate(string, number)`

## Description

Length argument must be enclosed in backticks; ex. "{{request.object.metadata.name | truncate(@, `9`)}}".

## Examples

```
truncate('foobar', `3`) == 'foo'
```

```
truncate('hello world', `5`) == 'hello'
```

```
truncate('short', `10`) == 'short'
```
