package lint

import (
	"fmt"

	yaml "k8s.io/apimachinery/pkg/util/yaml"
)

type FileFormatProcessor interface {
	ToJSON(input []byte) ([]byte, error)
}

type (
	JSONProcessor struct{}
	YAMLProcessor struct{}
)

func (p JSONProcessor) ToJSON(input []byte) ([]byte, error) {
	return input, nil
}

func (p YAMLProcessor) ToJSON(input []byte) ([]byte, error) {
	return yaml.ToJSON(input)
}

func getProcessor(format string, input []byte) (FileFormatProcessor, error) {
	if format == "" {
		format = detectFormat(input)
	}
	switch format {
	case ".json":
		return JSONProcessor{}, nil
	case ".yaml":
		return YAMLProcessor{}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}

func detectFormat(input []byte) string {
	if yaml.IsJSONBuffer(input) {
		return ".json"
	}
	return ".yaml"
}
