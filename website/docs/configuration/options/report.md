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
