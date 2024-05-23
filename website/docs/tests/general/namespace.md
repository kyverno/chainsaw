# Namespace management


Kubernetes organizes resources into two primary scopes: namespaced and cluster-scoped.

While namespaced resources belong to a specific namespace, cluster-scoped resources span across the entire Kubernetes cluster.

### Purpose of the Namespacer Interface

The [Namespacer interface](https://github.com/kyverno/chainsaw/blob/main/pkg/runner/namespacer/namespacer.go) ensures automated and consistent namespace assignment to Kubernetes resources.

- **Automated Namespacing**

    Automatically assign namespaces to resources that don't have one.

- **Ephemeral Namespaces**

    Handles temporary namespaces for specific tasks.
