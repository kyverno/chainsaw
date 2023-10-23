# IDE completion

While there is no currently available Kubernetes controller to handle the Chainsaw configuration files, the CRD definitions may be handy for users to leverage coding assistance in their favorite IDE.

## CRD definitions

CRD definitions are provided in our [GitHub repository](https://github.com/kyverno/chainsaw/tree/main/config/crds).

## Apply in cluster

If necessary, CRDs can be registered in your cluster with:

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/chainsaw/main/config/crds/chainsaw.kyverno.io_configurations.yaml
kubectl apply -f https://raw.githubusercontent.com/kyverno/chainsaw/main/config/crds/chainsaw.kyverno.io_tests.yaml
kubectl apply -f https://raw.githubusercontent.com/kyverno/chainsaw/main/config/crds/chainsaw.kyverno.io_teststeps.yaml
```
