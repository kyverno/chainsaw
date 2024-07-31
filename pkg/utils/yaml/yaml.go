package yaml

import (
	"gopkg.in/yaml.v3"
)

func remarshal(document []byte, unmarshal func(in []byte, out interface{}) (err error)) ([]byte, error) {
	if unmarshal == nil {
		unmarshal = yaml.Unmarshal
	}
	var pre map[string]any
	err := unmarshal(document, &pre)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(pre)
}

func Remarshal(document []byte) ([]byte, error) {
	return remarshal(document, nil)
}
