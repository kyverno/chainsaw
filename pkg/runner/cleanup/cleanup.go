package cleanup

import (
	"github.com/kyverno/chainsaw/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Cleaner = func(ctrlclient.Object, client.Client)
