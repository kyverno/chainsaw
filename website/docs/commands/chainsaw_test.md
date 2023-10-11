## chainsaw test

Stronger tool for e2e testing

```
chainsaw test [flags]... [test directories]...
```

### Options

```
      --config string          Chainsaw configuration file.
      --duration duration      The duration to use as default for configuration. (default 30s)
      --fullName               Use full test case folder path instead of folder name.
  -h, --help                   help for test
      --namespace string       Namespace to use for tests.
      --parallel int           The maximum number of tests to run at once. (default 8)
      --reportFormat string    Test report format (JSON|XML|nil).
      --reportName string      The name of the report to create. (default "kuttl-report")
      --skipDelete             If set, do not delete the resources after running the tests.
      --skipTestRegex string   Regular expression to skip tests based on.
      --stopOnFirstFailure     Stop the test upon encountering the first failure.
      --suppress strings       Logs to suppress.
      --testDirs strings       Directories containing test cases to run.
```

### SEE ALSO

* [chainsaw](chainsaw.md)	 - 

