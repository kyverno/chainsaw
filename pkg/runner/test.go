package runner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
)

func testName(config v1alpha1.ConfigurationSpec, test discovery.Test) (string, error) {
	if !config.FullName {
		return test.GetName(), nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(test.BasePath)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s[%s]", rel, test.GetName()), nil
}
