# Reports

Chainsaw can generate JUnit reports in `XML` or `JSON` format.

To produce a test report, configure the report format and report name in the configuration or using CLI flags.

## Configuration

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Configuration
metadata:
  name: custom-config
spec:
  # ...
  reportFormat: JSON
  reportName: chainsaw-report.json
  # ...
```

## Flag

```bash
chainsaw test --report-format JSON --report-name chainsaw-report.json ...
```
