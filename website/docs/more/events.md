# Working with events

Kubernetes events are regular Kubernetes objects and can be asserted on just like any other object:

```yaml
apiVersion: v1
kind: Event
reason: Started
source:
  component: kubelet
involvedObject:
  apiVersion: v1
  kind: Pod
  name: my-pod
```
