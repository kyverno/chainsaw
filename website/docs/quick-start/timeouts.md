# Control your timeouts

Timeouts in Chainsaw are specified per type of operation.
This is handy because the timeout varies greatly depending on the nature of an operation.

For example, applying a manifest in a cluster is expected to be reasonably fast, while validating a resource can be a long operation.

## Inheritance

Timeouts can be configured globally and at the test, step or individual operation level.

All timeouts configured at a given level are automatically inherited in child levels. When looking up a timeout, the most specific one takes precedence over the others.

!!! info
    To learn more about timeouts and how to configure global values, see the [timeouts configuration](../configuration/options/timeouts.md) page.

## At the test level

When a timeout is configured at the test level it will apply to all operations and steps in the test, unless overridden at a more specific level.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  # timeouts configured at the test level will apply to all operations and steps
  # unless overriden at the step level and/or individual operation level
  timeouts:
    apply: 5s
    assert: 1m
    # ...
  steps:
  - try:
    - apply:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          spec:
            storage:
              secret:
                name: minio
                type: s3
            # ...
    - assert:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          status:
            (conditions[?type == 'Ready']):
            - status: 'True'
```

## At the step level

When a timeout is configured at the step level it will apply to all operations in the step, unless overridden at a more specific level.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
    # timeouts configured at the step level will apply to all operations
    # in the step unless overriden at the individual operation level
  - timeouts:
      apply: 5s
      # ...
    try:
    - apply:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          spec:
            storage:
              secret:
                name: minio
                type: s3
            # ...
    # timeouts configured at the step level will apply to all operations
    # in the step unless overriden at the individual operation level
  - timeouts:
      assert: 1m
      # ...
    try:
    - assert:
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          status:
            (conditions[?type == 'Ready']):
            - status: 'True'
```

## At the operation level

When a timeout is configured at the operation level, it takes precedence over all timeouts configured at upper levels.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        # timeout configured at the operation level takes precedence
        # over timeouts configured at upper levels
        timeout: 5s
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          spec:
            storage:
              secret:
                name: minio
                type: s3
            # ...
    - assert:
        # timeout configured at the operation level takes precedence
        # over timeouts configured at upper levels
        timeout: 1m
        resource:
          apiVersion: tempo.grafana.com/v1alpha1
          kind: TempoStack
          metadata:
            name: simplest
          status:
            (conditions[?type == 'Ready']):
            - status: 'True'
```

## Next step

In the next section, we will see how Chainsaw [manages cleanup](./cleanup.md).
