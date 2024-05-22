# Reporting options

Report contains info about the report.

## Supported elements

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `format` | [`ReportFormatType`](#chainsaw-kyverno-io-v1alpha2-ReportFormatType) |  |  | <p>ReportFormat determines test report format (JSON|XML).</p> |
| `path` | `string` |  |  | <p>ReportPath defines the path.</p> |
| `name` | `string` |  |  | <p>ReportName defines the name of report to create. It defaults to "chainsaw-report".</p> |

## Configuration

### With file

```yaml
```

### With flags

```bash
```






Chainsaw can generate JUnit reports in `XML` or `JSON` format.

To produce a test report, configure the report format, report path and report name in the configuration or using CLI flags.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  # ...
  reportFormat: JSON
  reportName: chainsaw-report
  reportPath: /home/chainsaw
  # ...
```

## Flag

```bash
chainsaw test --report-format JSON --report-name chainsaw-report --report-path /path/to/save/report ...
```

> Note: The reportPath can be specified as either a relative or an absolute path.
