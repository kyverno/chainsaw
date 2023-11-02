package names

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
)

func Test(config v1alpha1.ConfigurationSpec, test discovery.Test) (string, error) {
	if test.Test == nil {
		return "", errors.New("test must not be nil")
	}
	if !config.FullName {
		return test.GetName(), nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working dir (%w)", err)
	}
	abs, err := filepath.Abs(test.BasePath)
	if err != nil {
		return "", fmt.Errorf("failed to compute absolute path for %s (%w)", test.BasePath, err)
	}
	rel, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path from %s to %s (%w)", cwd, abs, err)
	}
	return fmt.Sprintf("%s[%s]", rel, test.GetName()), nil
}
