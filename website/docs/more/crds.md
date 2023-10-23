# Custom resource definitions

New Custom Resource Definitions are not immediately available for use in the Kubernetes API until the Kubernetes API has acknowledged them.

If a Custom Resource Definition is being defined inside of a test step, be sure to to wait for the CustomResourceDefinition object to appear.

For example, given this Custom Resource Definition:

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.mycrd.k8s.io
spec:
  group: mycrd.k8s.io
  version: v1alpha1
  names:
    kind: MyCRD
    listKind: MyCRDList
    plural: mycrds
    singular: mycrd
  scope: Namespaced
```

Create the following assert:

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.mycrd.k8s.io
status:
  acceptedNames:
    kind: MyCRD
    listKind: MyCRDList
    plural: mycrds
    singular: mycrd
  storedVersions:
  - v1alpha1
```

And then the CRD can be used in subsequent steps:

```yaml
apiVersion: mycrd.k8s.io/v1alpha1
kind: MyCRD
metadata:
  name: my-crds
spec:
  test: test
```