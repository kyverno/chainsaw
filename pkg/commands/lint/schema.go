package lint

import (
	"fmt"
	"io/fs"

	"github.com/kyverno/chainsaw/pkg/data"
)

func getScheme(kind string, version string) ([]byte, error) {
	schemasFs, err := data.Schemas()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(schemasFs, fmt.Sprintf("%s-chainsaw-%s.json", kind, version))
}
