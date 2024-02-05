package values

import (
	"fmt"
	"io"
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

func Load(paths ...string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	for _, path := range paths {
		currentMap := map[string]interface{}{}
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

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
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
