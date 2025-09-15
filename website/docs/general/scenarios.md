# Scenarios

Sometimes you want to run the same test with different inputs. Scenarios can be used to define the different input sets the test will run with.

Input sets will be made available to the test through [bindings](./bindings.md).

## Syntax

The test below will run twice with different input bindings.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: scenarios
spec:
  # define two scenarios with different `message` binding
  scenarios:
  - bindings:
    - name: message
      value: hello
  - bindings:
    - name: message
      value: goodbye
  steps:
  - try:
    - script:
        env:
        - name: message
          value: ($message)
        content: echo $message
        # depending on the scenario ID, check the ouput of the script
        check:
          (trim_space($stdout)): ((to_number($test.scenarioId) == `1` && 'hello') || 'goodbye')
```

### Templating

Scenarios are treated early when executing a test.

The bindings created can then be used to template other test fields.

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: scenarios-bindings
spec:
  # define two scenarios to run the test against two different clusters
  scenarios:
  - bindings:
    - name: cluster
      value: ...
  - bindings:
    - name: cluster
      value: ...
  # use the bindings declared in each scenario to bind to the corresponding cluster
  cluster: ($cluster)
```

In the test above, the target cluster the test will run against will be driven by bindings declared in each scenario.
