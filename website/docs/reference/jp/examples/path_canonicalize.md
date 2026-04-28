# path_canonicalize

## Signature

`path_canonicalize(string)`

## Description

Normalizes or canonicalizes a given path by removing excess slashes.

## Examples

```
path_canonicalize('/foo//bar') == '/foo/bar'
```

```
path_canonicalize('/foo/./bar') == '/foo/bar'
```

```
path_canonicalize('/foo/../bar') == '/bar'
```
