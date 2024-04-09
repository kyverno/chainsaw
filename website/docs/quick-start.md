# Quick start

To use Chainsaw you will need a [Kubernetes](https://kybernetes.io) cluster, Chainsaw won't create one for you.

In this Quick start we will use [kind](https://kind.sigs.k8s.io) but feel free to use the tool of your choice.

!!! warning "Not a cluster management tool"

    We consider this is not the responsibility of Chainsaw to manage clusters.

    There are plenty of solutions to create and manage local clusters that will do that better than Chainsaw.

## Create a [kind](https://kind.sigs.k8s.io) cluster

Please refer to the [kind install docs](https://kind.sigs.k8s.io/docs/user/quick-start/#installation) to install it locally.

Once [kind](https://kind.sigs.k8s.io) creating a local cluster is as simple as running:

```bash
# define kind image
export KIND_IMAGE="kindest/node:v1.28.0"

# create cluster
kind create cluster --image $KIND_IMAGE
```

## Writing tests

A Chainsaw test is [made of YAML files in a folder](./tests/index.md).

YAML files can contain raw manifests with a special file naming convention to identify the step operations.
This is useful to create tests quickly but doesn't allow great flexibility.

Another option is to have a `chainsaw-test.yaml` file containing one (or more) `Test` resource. While more verbose, this offers full flexibility over the test and test steps configuration.

### Quick start

For this quick start, we will create a `Test` with one step and two operations:

1. Create a `ConfigMap` from a manifest
1. Verify the `ConfigMap` was created and contains the expected data

### Create the test folder

```bash
# create test folder
mkdir chainsaw-quick-start

# enter test folder
cd chainsaw-quick-start
```

### Create a `ConfigMap`

```bash
# create a ConfigMap
cat > configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: chainsaw-quick-start
data:
  foo: bar
EOF
```

### Create the test

By default, Chainsaw will look for a file named `chainsaw-test.yaml` in every folder.

```bash
# create test file
cat > chainsaw-test.yaml << EOF
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: quick-start
spec:
  steps:
  - try:
    # first operation: create the config map
    - apply:
        # file is relative to the test folder
        file: configmap.yaml
    # second operation: verify the config map exists and contains the expected data
    - assert:
        # file is relative to the test folder
        file: configmap.yaml
EOF
```

## Run Chainsaw

We finished writing our first test, now we can run Chainsaw to verify the test runs as expected:

```bash
chainsaw test

Version: (devel)
Loading default configuration...
- Using test file: chainsaw-test.yaml
- TestDirs [.]
- SkipDelete false
- FailFast false
- ReportFormat ''
- ReportName 'chainsaw-report'
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
### Chainsaw Charaterstics 

**Resource diff in assertion failures:** Chainsaw offers detailed resource diffs upon assertion failures. It provides a clear view of discrepancies between expected and actual resource.

**Resource templating support:** Chainsaw simplifies dynamic resource configuration with it's templating support. Instead on relying on `envsubst` for dynamic substitution of env-variable. This eliminates the need of any preprocessing step.

## Chainsaw Characteristics

**Simplified Testing of Kubernetes Operators**

Developing Kubernetes operators for complex systems, such as managing distributed database systems across multiple clusters, presents unique challenges. These operators need to dynamically respond to changing conditions, maintain system integrity, and ensure high availability and performance. Traditional testing methods often fall short in adequately simulating these dynamic environments or verifying the nuanced behaviors of these systems.

**Chainsaw: A Transformation in Testing**

Chainsaw emerges as a transformative solution, offering declarative, dynamic testing that aligns closely with Kubernetes' own principles. Here's how Chainsaw redefines testing workflows:

- **Declarative Testing with Dynamic Configurations**: Chainsaw enables teams to define test scenarios in YAML, incorporating dynamic configurations through resource templating. This approach eliminates the need for cumbersome scripting, enabling rapid setup of complex test environments.

Declarative Syntax:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    # ...
    - command:
        entrypoint: echo
        args:
        - hello chainsaw
    # ...
```
Resource Templating Support:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: template
spec:
  template: true
  steps:
  - assert:
      resource:
        # apiVersion, kind, name, namespace and labels are considered for templating
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: ($namespace)
        # other fields are not (they are part of the assertion tree)
        data:
          foo: ($namespace)
```
- **Advanced Assertions for Complex Logic**: With Chainsaw, teams can express complex logical conditions within their test definitions. Whether verifying the scaling behavior of a database cluster or ensuring data consistency across instances, Chainsaw's assertion model provides the precision and flexibility needed for comprehensive validation.

Logical Assertion Statements:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    # ...
    - assert:
        resource:
          apiVersion: v1
          kind: Deployment
          metadata:
            name: foo
          spec:
            (replicas > 3): true
    # ...
```

- **Streamlined Debugging and Maintenance**: Chainsaw's clear, concise test scenarios and detailed logging facilitate easier debugging and issue resolution. Teams can quickly identify problems and make adjustments, significantly reducing the maintenance overhead associated with custom scripts or rigid testing frameworks.

Resource diff in assertion failures:

```bash
    | 09:55:50 | deployment | step-1   | ASSERT    | RUN   | v1/Pod @ chainsaw-rare-liger/*
    | 09:56:20 | deployment | step-1   | ASSERT    | ERROR | v1/Pod @ chainsaw-rare-liger/*
        === ERROR
        ---------------------------------------------------
        v1/Pod/chainsaw-rare-liger/example-5477b4ff8c-tnhd9
        ---------------------------------------------------
        * metadata.annotations.foo: Invalid value: "null": Expected value: "bar"
        
        --- expected
        +++ actual
        @@ -1,10 +1,16 @@
         apiVersion: v1
         kind: Pod
         metadata:
        -  annotations:
        -    foo: bar
           labels:
             app: nginx
        +  name: example-5477b4ff8c-tnhd9
           namespace: chainsaw-rare-liger
        +  ownerReferences:
        +  - apiVersion: apps/v1
        +    blockOwnerDeletion: true
        +    controller: true
        +    kind: ReplicaSet
        +    name: example-5477b4ff8c
        +    uid: 118abe16-ec42-4894-83db-64479c4aac6f
         spec: {}
    | 09:56:20 | deployment | step-1   | TRY       | DONE  |

```

**Example Test Scenario: Verifying Operator Scaling Behavior**

Consider a scenario where a Kubernetes operator manages a distributed database. The operator needs to scale database instances based on load, ensuring optimal performance and resource utilization. Here's a simplified Chainsaw test that verifies this functionality:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: database-scaling
spec:
  steps:
  - try:
      apply: db-cluster.yaml
      values:
        replicaCount: 3
  - assert:
      resource:
        apiVersion: v1
        kind: DatabaseCluster
        metadata:
          name: test-cluster
      status:
        replicas: 3
  - try:
      patch:
        file: db-cluster-scale.yaml
        values:
          newReplicaCount: 5
  - assert:
      resource:
        apiVersion: v1
        kind: DatabaseCluster
        metadata:
          name: test-cluster
      status:
        (replicas > 3): true
```

This test dynamically configures a database cluster, verifies its initial scale, applies a scaling operation, and asserts the successful scaling based on the specified conditions.