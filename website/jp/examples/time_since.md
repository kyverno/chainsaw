!!! note "Non-deterministic output"

    When the third argument is empty, `time_since` uses the current time as the end point and produces a different value each call.

```
# calculate elapsed time between two fixed points
time_since('', '2024-01-01T00:00:00Z', '2024-01-02T06:00:00Z') == '30h0m0s'
```

```
# calculate elapsed time with a custom layout
time_since('2006-01-02', '2024-01-01', '2024-01-15') == '336h0m0s'
```
