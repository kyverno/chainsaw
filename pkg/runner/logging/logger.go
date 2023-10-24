package logging

import (
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
}

func ResourceOp(logger Logger, op string, key ctrlclient.ObjectKey, obj ctrlclient.Object, attempts int, err error) {
	if logger == nil {
		return
	}
	var options []string
	if err != nil {
		options = append(options, fmt.Sprintf("err=%s", err))
	}
	if attempts != 0 {
		options = append(options, fmt.Sprintf("att=%d", attempts))
	}
	gvk := obj.GetObjectKind().GroupVersionKind()
	if len(options) != 0 {
		logger.Logf("%s[%s/%s] %s (%s)", strings.ToUpper(op), gvk.GroupVersion(), gvk.Kind, client.Name(key), strings.Join(options, ","))
	} else {
		logger.Logf("%s[%s/%s] %s", strings.ToUpper(op), gvk.GroupVersion(), gvk.Kind, client.Name(key))
	}
}
