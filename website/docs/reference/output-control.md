# Output Control

Chainsaw provides several options to control the output format and verbosity when running tests.

## Suppressing Warnings

Chainsaw tests may generate warning messages that are helpful during development but can clutter output in CI/CD environments or when warnings are expected. The `--no-warnings` flag allows you to suppress these warning messages.

### Usage

```bash
chainsaw test --no-warnings [other flags] [test directories]
```

### Example

When running tests with expected warnings:

```bash
# Without suppression (shows warnings)
chainsaw test ./tests

# With warning suppression
chainsaw test --no-warnings ./tests
```

### When to Use

The `--no-warnings` flag is particularly useful in the following scenarios:

1. **CI/CD Pipelines**: Keeping build logs clean and focused on relevant information
2. **Tests with Expected Warnings**: When your tests intentionally trigger conditions that generate warnings
3. **Large Test Suites**: Reducing output verbosity when running many tests

### Output Comparison

#### With Warnings (Default)
```
| 15:10:23 | test-name | step-1 | APPLY     | OK    |
| 15:10:24 | test-name | step-2 | CREATE    | WARN  | Timeout value is low
| 15:10:25 | test-name | step-3 | DELETE    | OK    |
```

#### With Warnings Suppressed
```
| 15:10:23 | test-name | step-1 | APPLY     | OK    |
| 15:10:25 | test-name | step-3 | DELETE    | OK    |
```

## Other Output Control Options

In addition to suppressing warnings, Chainsaw offers other output control options:

- **--no-color**: Disables colored output, useful for terminals that don't support ANSI colors or when redirecting output to files
- **--report-format**: Changes the format of test reports (JSON, XML, etc.)

## Combining Options

Output control options can be combined for further customization:

```bash
chainsaw test --no-warnings --no-color ./tests
```

This gives you clean, uncolored output that's ideal for automated environments and log processing. 