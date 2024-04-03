package lint

import (
	"fmt"
	"io/fs"
	"path"

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
	return fs.ReadFile(data.Schemas(), path.Join("schemas", "json", "test-chainsaw-v1alpha1.json"))
}

func getConfigurationSchema() ([]byte, error) {
	return fs.ReadFile(data.Schemas(), path.Join("schemas", "json", "configuration-chainsaw-v1alpha1.json"))
}
