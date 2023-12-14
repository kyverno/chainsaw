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

func getProcessor(format string) (FileFormatProcessor, error) {
	if format == "" {
		return nil, fmt.Errorf("file format must be specified")
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
