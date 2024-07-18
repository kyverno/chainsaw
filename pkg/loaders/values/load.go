package values

import (
	"fmt"
	"io"
	"os"
	"strings"

	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"sigs.k8s.io/yaml"
)

func Load(paths ...string) (map[string]any, error) {
	base := map[string]any{}
	for _, path := range paths {
		currentMap := map[string]any{}
		bytes, err := readFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, fmt.Errorf("failed to parse %s (%w)", path, err)
		}
		base = mapsutils.Merge(base, currentMap)
	}
	return base, nil
}

func readFile(filePath string) ([]byte, error) {
	if strings.TrimSpace(filePath) == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filePath)
}
