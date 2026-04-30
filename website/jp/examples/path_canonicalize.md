```
path_canonicalize('/foo//bar') == '/foo/bar'
```

```
path_canonicalize('/foo/./bar') == '/foo/bar'
```

```
path_canonicalize('/foo/../bar') == '/bar'
```
