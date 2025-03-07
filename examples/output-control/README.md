# Output Control Example

This example demonstrates how to use Chainsaw's output control features, specifically the `--no-warnings` flag.

## Overview

The test in this directory intentionally generates warning messages that can be suppressed using the `--no-warnings` flag.

## Files

- `chainsaw-test.yaml`: Contains a test that generates a warning message

## Usage

Run the test with and without the `--no-warnings` flag to see the difference:

```bash
# Run with warnings (default)
chainsaw test ./examples/output-control

# Run with warnings suppressed
chainsaw test --no-warnings ./examples/output-control
```

## Expected Output

### With Warnings (Default)

```
| XX:XX:XX | warning-example | step-warning-example | COMMAND   | WARN  | This is a warning message that would be suppressed with --no-warnings
| XX:XX:XX | warning-example | step-warning-example | COMMAND   | DONE  |
```

### With Warnings Suppressed

```
| XX:XX:XX | warning-example | step-warning-example | COMMAND   | DONE  |
```

The warning message is completely removed from the output when using the `--no-warnings` flag.

## Additional Resources

For more information on output control options, see the [Output Control documentation](https://kyverno.github.io/chainsaw/latest/reference/output-control/). 