# Reports

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
