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
