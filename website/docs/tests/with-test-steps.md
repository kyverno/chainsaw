# With test steps

Chainsaw has a granular approach to testing using individual `TestStepSpec` elements. This document focuses on understanding test steps.

## The `TestStepSpec` Component

Each `TestStepSpec` represents a specific stage in your test and offers detailed control over the process. Here's what you can define within each test step:

- **Timeout**: Dictates how long the test step should run before being marked as failed due to a timeout.
- **Assert**: Specifies the conditions that must be true for the step to pass. Essentially, it's where you set your expectations.
- **Apply**: Denotes the Kubernetes resources or configurations that should be applied at this stage.
- **Error**: Lists the expected errors for this step. This is vital for cases where certain errors are anticipated and should be treated as part of the expected behavior.
- **Delete**: Points out resources that need to be removed before this step gets executed. It ensures the desired state of the environment before the step runs.

## Crafting a Test Step in YAML

Here's an example of a test step that applies a resource and asserts its state:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: TestStep
metadata:
  name: test-step-name
steps:
# Step to apply a resource.
  apply:
  - file: my-resource.yaml  # Apply a specific resource.

# Step to assert the applied resource's state.
  assert:
  - file: expected-resource-state.yaml  # Compare against an expected state.

# Step to handle expected errors.
  error:
  - file: expected-error.yaml  # Compare against an expected error.

# Step to delete resources before the next action.
  delete:
  - file: resource-to-delete.yaml
