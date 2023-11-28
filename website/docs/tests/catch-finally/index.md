# What is catch / finally

Catch and Finally are additional fields to collect certain information about the outcome of a step should it fail (in the case of `catch`) or at the end of the step (in the case of `finally`).

The ultimate goal of collectors is to gather information about the failure of a step and therefore help understand what caused it to fail.

A test step can have an arbitrary number of collectors. A collector can be configured with a `timeout`.

!!! note

    `catch` operations are only invoked in cases where **a failure occurs** and not if the step succeeds.

    `finally` operations are always invoked regardles of the test success or failure.

Collection can occur from:

- Pod logs
- Namespace events
- Or the output of a custom command or script

## Catch lifecycle

The `catch` operations are executed **after** the step failed, and **before** the step cleanup happens.

!!! info "Catch lifecycle"

    1. The step starts executing
    1. An operation fails (**before catch operations are executed**)
    1. Catch operations are executed
    1. The step cleanup executes (**after catch operations are executed**)

This is important that collectors run before cleanup so that they have a chance to collect logs from pods responsible for the step failure.

## Configuration

Catch / Finally are a per step configuration and are registered under the `catch` / `finally` sections of a test step spec.

!!! example "Collect pod logs"

    ```yaml
        try:
        # ...
        catch:
        - podLogs:
            name: my-pod
        finally:
        - podLogs:
            name: my-pod
    ```
    See [Pod logs](pod-logs.md) for details and supported configurations.

!!! example "Collect events"

    ```yaml
        try:
        # ...
        catch:
        - events: {}
        finally:
        - events: {}
    ```
    See [Events](events.md) for details and supported configurations.

!!! example "Execute a custom command"

    ```yaml
        try:
        # ...
        catch:
        - command:
            entrypoint: kubectl
            args:
            - get
            - pod
            - -n
            - $NAMESPACE
        finally:
        - command:
            entrypoint: kubectl
            args:
            - get
            - pod
            - -n
            - $NAMESPACE
    ```
    See [Commands](commands.md) for details and supported configurations.

!!! example "Execute a custom script"

    ```yaml
        try:
        # ...
        catch:
        - script:
            content: |
              echo "an error has occured"
        finally:
        - script:
            content: |
              echo "goodbye"
    ```
    See [Scripts](scripts.md) for details and supported configurations.

!!! example "Execute a custom script with timeout"

    ```yaml
        try:
        # ...
        catch:
        - script:
            content: |
              echo "an error has occured"
        finally:
        - script:
            content: |
              echo "goodbye"
          timeout: 15s
    ```
    See [Scripts](scripts.md) for details and supported configurations.
