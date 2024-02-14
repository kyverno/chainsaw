# Timeouts

Timeouts in Chainsaw are specified per type of operation.
This is required because the timeout varies greatly depending on the nature of an operation.

For example, applying a manifest in a cluster is expected to be reasonably fast, while validating a resource can be a long operation.

Chainsaw supports separately configuring the timeouts below:

- **Apply**

    When Chainsaw applies manifests in a cluster

- **Assert**

    When Chainsaw validates resources in a cluster

- **Cleanup**

    When Chainsaw removes resources from a cluster created for a test

- **Delete**

    When Chainsaw deletes resources from a cluster

- **Error**

    When Chainsaw validates resources in a cluster

- **Exec**

    When Chainsaw executes arbitrary commands or scripts

!!! note "Overriding timeouts"

    Each timeout can be overridden at the test level, test step level, or individual operation level.

    Timeouts defined in the `Configuration` are used in operations when not overridden.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: custom-config
spec:
  # ...
  timeouts:
    apply: 45s
    assert: 20s
    cleanup: 45s
    delete: 25s
    error: 10s
    exec: 45s
  # ...
```

## Flag

```bash
chainsaw test                     \
    --apply-timeout 45s             \
    --assert-timeout 45s            \
    --cleanup-timeout 45s           \
    --delete-timeout 45s            \
    --error-timeout 45s             \
    --exec-timeout 45s              \
    ...
```
