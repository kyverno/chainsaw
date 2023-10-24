# Chainsaw Tool Configuration

Chainsaw is a comprehensive tool designed to facilitate end-to-end (e2e) testing in Kubernetes. This documentation will focus on providing a breakdown of its configuration structure and how to use it.

## Configuration Spec (v1alpha1)

## 1. Overview:

Chainsaw is described as a "Stronger tool for e2e testing". With its versatile configuration options, you can customize the testing process to fit your needs.

## 2. Configuration Loading Process:

Chainsaw prioritizes its configuration in the following order:

1. **User-Specified Configuration**: If you explicitly provide a configuration file using a command-line flag.
2. **Default Configuration File**: If no configuration is specified, Chainsaw will look for a default file named `.chainsaw.yaml` in the current directory.
3. **Internal Default Configuration**: In the absence of both the above, Chainsaw will use its built-in default configuration.

### How to Specify a Configuration:

To use a custom configuration file:

```bash
chainsaw test --config path/to/your/config.yaml
```

If you don't specify any configuration, Chainsaw will look for the default configuration file in the current directory. If that's not found, it will fall back to its internal default settings.

### Overriding with Flags:

Even after a configuration is loaded, you can override specific settings using command-line flags:

```bash
chainsaw test --timeout 45s
```

In this example, the timeout configuration will be set to 45 seconds, regardless of the value in the loaded configuration file.

### Fields

| **Field**            | **Description**                                              | **Default**                   |
|----------------------|--------------------------------------------------------------|-------------------------------|
| Timeout              | Global timeout for tests.                                    | 30 seconds                    |
| TestDirs             | Directories containing test cases.                           |                               |
| SkipDelete           | Prevents resource deletion post-testing.                     | false                         |
| FailFast             | Stops test execution on first failure.                       | false                         |
| Parallel             | Max number of simultaneous tests.                            | 8                             |
| RepeatCount          | How many times tests should be executed.                     | 1                             |
| ReportFormat         | Test report format (`JSON`, `XML`, or none).                 |                               |
| ReportName           | Name of the report.                                          | "chainsaw-report"             |
| Namespace            | Namespace for tests.                                         |                               |
| Suppress             | Suppress specific logs.                                      |                               |
| FullName             | Use full test case folder path for representation.           | false                         |
| ExcludeTestRegex     | Exclude tests using regex.                                   |                               |
| IncludeTestRegex     | Include tests using regex.                                   |                               |

---

### Flags

| **Flag**                   | **Description**                                              | **Default**            |
|----------------------------|--------------------------------------------------------------|------------------------|
| `--timeout`                | Default timeout for the configuration.                        | 30 seconds             |
| `--config`                 | Chainsaw configuration file.                                  | -                      |
| `--test-dir`               | Directories containing test cases to run.                     | -                      |
| `--skip-delete`            | Prevents resource deletion post-testing.                     | -                      |
| `--stop-on-first-failure`  | Stops test execution on first failure.                       | -                      |
| `--parallel`               | Max number of tests running concurrently.                    | 8                      |
| `--repeat-count`           | Number of times to repeat each test.                          | -                      |
| `--report-format`          | Test report format (`JSON`, `XML`, or none).                 | -                      |
| `--report-name`            | Name for the test report.                                    | "chainsaw-report"      |
| `--namespace`              | Namespace for tests.                                         | -                      |
| `--suppress`               | Lists the logs to suppress.                                  | -                      |
| `--full-name`              | Use the full test case folder path instead of folder name.   | -                      |
| `--include-test-regex`     | Regex to include specific tests.                             | -                      |
| `--exclude-test-regex`     | Regex to exclude specific tests.                             | -                      |

### Usage Example

```bash
chainsaw test --config my-config.yaml --test-dir /path/to/tests --parallel 10
```

This command will run tests using the configuration from my-config.yaml, taking tests from /path/to/tests, and running a maximum of 10 tests simultaneously.
