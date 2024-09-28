package kubectl

import (
	"context"
	"errors"
	"path"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Proxy(ctx context.Context, client client.Client, tc apis.Bindings, collector *v1alpha1.Proxy) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := collector.Name.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := collector.Namespace.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	targetPath, err := collector.TargetPath.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	targetPort, err := collector.TargetPort.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	resource, _, err := mapResource(ctx, client, tc, collector.ObjectType)
	if err != nil {
		return "", nil, err
	}
	args := []string{"get", "--raw", path.Join("/api", "v1", "namespaces", namespace, resource, name+":"+targetPort, "proxy", targetPath)}
	return "kubectl", args, nil
}
