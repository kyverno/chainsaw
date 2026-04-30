```
label_match(`{"app":"foo","env":"prod"}`, `{"app":"foo"}`) == `true`
```

```
label_match(`{"app":"foo","env":"prod"}`, `{"app":"foo","env":"prod"}`) == `true`
```

```
label_match(`{"app":"foo"}`, `{"app":"foo","env":"prod"}`) == `false`
```

```
label_match(`{"app":"foo"}`, `{}`) == `true`
```
