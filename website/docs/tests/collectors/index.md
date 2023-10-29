# What is a collector

Collectors are used to collect certain information about the outcome of a step should it fail.
The ultimate goal of collectors is to gather information about the failure of a step and therefore help understand what caused it to fail.

A test step can have an arbitrary number of collectors.

!!! note

    A collector is only invoked in cases where **a failure occurs** and not if the step succeeds.

Collection can occur from:

- Pod logs
- Namespace events
- Or the output of a custom command or script

## Collectors lifecycle

The collectors are executed **after** the step failed, and **before** the step cleanup happens.

!!! info "Collectors lifecycle"

    1. The step starts executing
    1. An operation fails (**before collectors are executed**)
    1. Collectors are executed
    1. The step cleanup executes (**after collectors are executed**)

This is important that collectors run before cleanup so that they have a chance to collect logs from pods responsible for the step failure.

## Configuration

Collectors are a per step configuration and are registered under the `onFailure` section of a test step spec.

!!! example "Collect pod logs"

    ```yaml
        onFailure:
        - collect:
            podLogs:
              name: my-pod
    ```
    See [Pod logs](pod-logs.md) for details and supported configurations.

!!! example "Collect events"

    ```yaml
        onFailure:
        - collect:
            events: {}
    ```
    See [Events](events.md) for details and supported configurations.

!!! example "Execute a custom command"

    ```yaml
        - exec:
            command:
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
        - exec:
            script:
              content: |
                echo "an error has occured"
    ```
    See [Scripts](scripts.md) for details and supported configurations.
