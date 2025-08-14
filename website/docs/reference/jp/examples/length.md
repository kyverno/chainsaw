# length

## Signature

`length(string|array|object)`

## Description

Returns the length of the given argument. If the argument is a string this function returns the number of code points in the string. If the argument is an array this function returns the number of elements in the array. If the argument is an object this function returns the number of key-value pairs in the object.

## Examples

```
length([`10`,`15`,`20`]) == `3`
```

```
length([]) == `0`
```

```
length(null) -> error
```
