# Run tests

After [installing chainsaw](./install.md) and [writing tests](./first-test.md), the next natural step is to run Chainsaw to execute the tests.

## Create a local cluster

To use Chainsaw you will need a [Kubernetes](https://kubernetes.io) cluster, **Chainsaw won't create one for you**.

!!! info "Not a cluster management tool"
    We consider this is not the responsibility of Chainsaw to manage clusters.
    There are plenty of solutions to create and manage local clusters that will do that better than Chainsaw.

The command below will create a local cluster using [kind](https://kind.sigs.k8s.io).
Use the tool of your choice or directly jump to the next section if you already have a `KUBECONFIG` configured and pointing to a valid cluster.

```bash
# create cluster
kind create cluster --image "kindest/node:v1.29.4"
```

## Run Chainsaw

Now you can run the `chainsaw test` command.

```
> chainsaw test

Version: (devel)
Loading default configuration...
- Using test file: chainsaw-test.yaml
- TestDirs [.]
- SkipDelete false
- FailFast false
- ReportFormat ''
- ReportName ''
- Namespace ''
- FullName false
- IncludeTestRegex ''
- ExcludeTestRegex ''
- ApplyTimeout 5s
- AssertTimeout 30s
- CleanupTimeout 30s
- DeleteTimeout 15s
- ErrorTimeout 30s
- ExecTimeout 5s
Loading tests...
- quick-start (.)
Running tests...
=== RUN   chainsaw
=== PAUSE chainsaw
=== CONT  chainsaw
=== RUN   chainsaw/quick-start
=== PAUSE chainsaw/quick-start
=== CONT  chainsaw/quick-start
    | 10:44:26 | quick-start | @setup   | CREATE    | OK    | v1/Namespace @ chainsaw-immense-jay
    | 10:44:26 | quick-start | step-1   | TRY       | RUN   |
    | 10:44:26 | quick-start | step-1   | APPLY     | RUN   | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | step-1   | CREATE    | OK    | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | step-1   | APPLY     | DONE  | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | step-1   | ASSERT    | RUN   | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | step-1   | ASSERT    | DONE  | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | step-1   | TRY       | DONE  |
    | 10:44:26 | quick-start | @cleanup | DELETE    | RUN   | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | @cleanup | DELETE    | OK    | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | @cleanup | DELETE    | DONE  | v1/ConfigMap @ chainsaw-immense-jay/chainsaw-quick-start
    | 10:44:26 | quick-start | @cleanup | DELETE    | RUN   | v1/Namespace @ chainsaw-immense-jay
    | 10:44:26 | quick-start | @cleanup | DELETE    | OK    | v1/Namespace @ chainsaw-immense-jay
    | 10:44:31 | quick-start | @cleanup | DELETE    | DONE  | v1/Namespace @ chainsaw-immense-jay
--- PASS: chainsaw (0.00s)
    --- PASS: chainsaw/quick-start (5.25s)
PASS
Tests Summary...
- Passed  tests 1
- Failed  tests 0
- Skipped tests 0
Done.
```

!!! tip
    Chainsaw expects a path to the test folder and will discover tests by analyzing files recursively. When no path is provided Chainsaw will use the current path by default (`.`).

## Next step

The test above demonstrates the most basic usage of Chainsaw. In the next sections, we will look at the [main features that make Chainsaw a very unique tool](./assertion-trees.md).
