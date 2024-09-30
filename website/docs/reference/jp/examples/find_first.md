# find_first

## Signature

`find_first(string, string, number, number)`

## Description

Returns the zero-based index of the first occurence where the substring appears in a string or null if it does not appear.

## Examples

```
find_first('subject string', 'string', `0`) == `8`
```

```
find_first('subject string', 'string', `8`) == `8`
```

```
find_first('subject string', 'string', `9`) == null
```

```
find_first('subject string', 'string', `0`, `14`) == `8`
```

```
find_first('subject string', 'string', `0`, `13`) == null
```
