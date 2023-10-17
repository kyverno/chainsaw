package logging

import (
	"strings"

	"github.com/kyverno/chainsaw/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
}

func ResourceOp(logger Logger, op string, key ctrlclient.ObjectKey, obj ctrlclient.Object) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	logger.Logf("%s[%s/%s] %s", strings.ToUpper(op), gvk.GroupVersion(), gvk.Kind, client.Name(key))
}
