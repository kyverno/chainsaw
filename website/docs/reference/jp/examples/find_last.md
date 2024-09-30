# find_last

## Signature

`find_last(string, string, number, number)`

## Description

Returns the zero-based index of the last occurence where the substring appears in a string or null if it does not appear.

## Examples

```
find_last('subject string', 'string', `0`) == `8`
```

```
find_last('subject string', 'string', `8`) == `8`
```

```
find_last('subject string', 'string', `9`) == null
```

```
find_last('subject string', 'string', `0`, `14`) == `8`
```

```
find_last('subject string', 'string', `0`, `13`) == null
```
