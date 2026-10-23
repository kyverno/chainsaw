# Reporting options

Reporting options contain the configuration used by Chainsaw for reporting.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `format` | `JSON` | ReportFormat determines test report format (JSON, XML, JUNIT-TEST, JUNIT-STEP, JUNIT-OPERATION). |
| `path` | | ReportPath defines the path. |
| `name` | `chainsaw-report` | ReportName defines the name of report to create. It defaults to "chainsaw-report". |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  report:
    format: JSON
    name: chainsaw-report
    path: /home/chainsaw
```

### With flags

!!! note
    The report path can be specified as either a relative or an absolute path.

```bash
chainsaw test                             \
  --report-format JSON                    \
  --report-name chainsaw-report           \
  --report-path /path/to/save/report
```

## Operation names and output

The report gives each operation of a step the `name` of the operation.
When an operation has no `name`, the report uses its `description`.
When an operation has neither, the report uses `operation <number>`.
When an operation has more than one resource, the name ends with the number of the resource, for example `apply-resources (2/3)`.

The report also keeps the log of each operation, without colors:

- In the `JUNIT-OPERATION` format, the log of an operation is in the `system-out` element of its test case.
- In the `JUNIT-STEP` format, the logs of all operations of a step are in the `system-out` element of the step test case.
- In the `JSON` format, the log of an operation is in its `output` field.

!!! note
    Only the operations of `try` and the deletes of the resources that a step created go to the report.
    The operations of `catch`, `finally` and `cleanup` do not go to the report.
