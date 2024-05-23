# Work with CRDs

New CRDs are not immediately available for use in the Kubernetes API until the Kubernetes API has acknowledged them.

If a CRD is being defined inside of a test step, be sure to wait for it to appear.

The test below applies a CRD and waits for it to become available:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  steps:
  - try:
    - apply:
        resource:
          apiVersion: apiextensions.k8s.io/v1
          kind: CustomResourceDefinition
          metadata:
            name: issues.example.com
          spec:
            group: example.com
            names:
              kind: Issue
              listKind: IssueList
              plural: issues
              singular: issue
            scope: Namespaced
            versions: ...
    - assert:
        resource:
          apiVersion: apiextensions.k8s.io/v1
          kind: CustomResourceDefinition
          metadata:
            name: issues.example.com
          status:
            acceptedNames:
              kind: Issue
              listKind: IssueList
              plural: issues
              singular: issue
            storedVersions:
            - v1alpha1
```

The CRD can be used in subsequent steps.
