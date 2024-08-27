package simple

import (
	"github.com/kyverno/chainsaw/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(cfg *rest.Config) (client.Client, error) {
	var opts ctrlclient.Options
	client, err := ctrlclient.New(cfg, opts)
	if err != nil {
		return nil, err
	}
	return ctrlclient.WithFieldValidation(client, metav1.FieldValidationStrict), nil
}
