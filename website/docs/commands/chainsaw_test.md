## chainsaw test

Stronger tool for e2e testing

```
chainsaw test [flags]... [test directories]...
```

### Options

```
      --config string                       Chainsaw configuration file.
      --exclude-test-regex string           Regular expression to exclude tests.
      --full-name                           Use full test case folder path instead of folder name.
  -h, --help                                help for test
      --include-test-regex string           Regular expression to include tests.
      --kube-as string                      Username to impersonate for the operation
      --kube-as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --kube-as-uid string                  UID to impersonate for the operation
      --kube-certificate-authority string   Path to a cert file for the certificate authority
      --kube-client-certificate string      Path to a client certificate file for TLS
      --kube-client-key string              Path to a client key file for TLS
      --kube-cluster string                 The name of the kubeconfig cluster to use
      --kube-context string                 The name of the kubeconfig context to use
      --kube-disable-compression            If true, opt-out of response compression for all requests to the server
      --kube-insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
  -n, --kube-namespace string               If present, the namespace scope for this CLI request
      --kube-password string                Password for basic authentication to the API server
      --kube-proxy-url string               If provided, this URL will be used to connect via proxy
      --kube-request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
      --kube-server string                  The address and port of the Kubernetes API server
      --kube-tls-server-name string         If provided, this name will be used to validate server certificate. If this is not provided, hostname used to contact the server is used.
      --kube-token string                   Bearer token for authentication to the API server
      --kube-user string                    The name of the kubeconfig user to use
      --kube-username string                Username for basic authentication to the API server
      --namespace string                    Namespace to use for tests.
      --parallel int                        The maximum number of tests to run at once. (default 8)
      --repeat-count int                    Number of times to repeat each test. (default 1)
      --report-format string                Test report format (JSON|XML|nil).
      --report-name string                  The name of the report to create. (default "chainsaw-report")
      --skip-delete                         If set, do not delete the resources after running the tests.
      --stop-on-first-failure               Stop the test upon encountering the first failure.
      --suppress stringArray                Logs to suppress.
      --test-dir stringArray                Directories containing test cases to run.
      --timeout duration                    The timeout to use as default for configuration. (default 30s)
```

### SEE ALSO

* [chainsaw](chainsaw.md)	 - 

