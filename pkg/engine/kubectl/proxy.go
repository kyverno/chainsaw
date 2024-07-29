package kubectl

import (
	"errors"
	"path"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
)

func Proxy(client client.Client, tc binding.Bindings, collector *v1alpha1.Proxy) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := bindings.String(collector.Name, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := bindings.String(collector.Namespace, tc)
	if err != nil {
		return "", nil, err
	}
	targetPath, err := bindings.String(collector.TargetPath, tc)
	if err != nil {
		return "", nil, err
	}
	targetPort, err := bindings.String(collector.TargetPort, tc)
	if err != nil {
		return "", nil, err
	}
	resource, _, err := mapResource(client, tc, collector.ObjectType)
	if err != nil {
		return "", nil, err
	}
	args := []string{"get", "--raw", path.Join("/api", "v1", "namespaces", namespace, resource, name+":"+targetPort, "proxy", targetPath)}
	return "kubectl", args, nil
}
