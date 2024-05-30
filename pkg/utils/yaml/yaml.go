package yaml

import (
	"gopkg.in/yaml.v3"
)

func Remarshal(document []byte) ([]byte, error) {
	var pre map[string]any
	err := yaml.Unmarshal(document, &pre)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(pre)
}
