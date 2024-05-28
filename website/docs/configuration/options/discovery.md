# Discovery options

Discovery options contain the discovery configuration used by Chainsaw when discovering tests in specified folders.

## Supported elements

| Element | Default | Description |
|---|---|---|
| `testFile` | `chainsaw-test` | TestFile is the name of the file containing the test to run. If no extension is provided, chainsaw will try with .yaml first and .yml if needed. |
| `fullName` | `false` | FullName makes use of the full test case folder path instead of the folder name. |
| `includeTestRegex` |  | IncludeTestRegex is used to include tests based on a regular expression. |
| `excludeTestRegex` |  | ExcludeTestRegex is used to exclude tests based on a regular expression. |

## Configuration

### With file

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: example
spec:
  discovery:
    testFile: chainsaw-test
    fullName: true
    includeTestRegex: chainsaw/.*
    excludeTestRegex: chainsaw/exclude-.*
```

### With flags

```bash
chainsaw test                                   \
  --test-file chainsaw-test                     \
  --full-name                                   \
  --include-test-regex 'chainsaw/.*'            \
  --exclude-test-regex 'chainsaw/exclude-.*'
```
