# Check Kubernetes version

The test below fetches the Kubernetes cluster version using `x_k8s_server_version`.
It then uses the minor version retrieved to adapt an assertion based on the value in the `$minorversion` binding.

!!!tip
    You can implement a ternary operator in JMESPath using an expression like this:
    
    `<condition> && <value-if-true> || <value-if-false>`

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: example
spec:
  bindings:
  - name: version
    value: (x_k8s_server_version($config))
  - name: minorversion
    value: (to_number($version.minor))
  steps:
  - try:
    - apply:
        resource:
          apiVersion: v1
          kind: Pod
          metadata:
            name: pod01
          spec:
            containers:
            - name: busybox
              image: busybox:1.35
    # ...
    - assert:
        resource:
          apiVersion: v1
          kind: Pod
          metadata:
            annotations:
              # If the minor version of the Kubernetes cluster against which this
              # is tested is less than 29, the annotation is expected to have the group 'system:masters' in it.
              # Otherwise, due to a change in kubeadm, the group should be 'kubeadm:cluster-admins'.
              kyverno.io/created-by: (($minorversion < `29` && '{"groups":["system:masters","system:authenticated"],"username":"kubernetes-admin"}') || '{"groups":["kubeadm:cluster-admins","system:authenticated"],"username":"kubernetes-admin"}')
            name: pod01
```