package values

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		base = mergeMaps(base, currentMap)
	}
	return base, nil
}

func mergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]any); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func readFile(filePath string) ([]byte, error) {
	if strings.TrimSpace(filePath) == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filePath)
}
