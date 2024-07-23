package lint

import (
	"fmt"
	"io/fs"

	"github.com/kyverno/chainsaw/pkg/data"
)

func getScheme(schema string) ([]byte, error) {
	switch schema {
	case "test":
		return getTestSchema()
	case "configuration":
		return getConfigurationSchema()
	default:
		return nil, fmt.Errorf("unknown schema: %s", schema)
	}
}

func getTestSchema() ([]byte, error) {
	schemasFs, err := data.Schemas()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(schemasFs, "test-chainsaw-v1alpha1.json")
}

func getConfigurationSchema() ([]byte, error) {
	schemasFs, err := data.Schemas()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(schemasFs, "configuration-chainsaw-v1alpha1.json")
}
