package names

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
)

type (
	workignDirInterface   = func() (string, error)
	absolutePathInterface = func(string) (string, error)
	relativePathInterface = func(string, string) (string, error)
)

func Test(test discovery.Test, full bool) (string, error) {
	if test.Test == nil {
		return "", errors.New("test must not be nil")
	}
	if !full {
		return test.Test.GetName(), nil
	}
	return helpTest(test, nil, nil, nil)
}

func Step(step v1alpha1.TestStep, i int) string {
	if step.Name != "" {
		return step.Name
	}
	return fmt.Sprintf("step-%d", i+1)
}

func helpTest(test discovery.Test, workingDir workignDirInterface, absolutePath absolutePathInterface, relativePath relativePathInterface) (string, error) {
	if workingDir == nil {
		workingDir = os.Getwd
	}
	if absolutePath == nil {
		absolutePath = filepath.Abs
	}
	if relativePath == nil {
		relativePath = filepath.Rel
	}
	cwd, err := workingDir()
	if err != nil {
		return "", fmt.Errorf("failed to get current working dir (%w)", err)
	}
	abs, err := absolutePath(test.BasePath)
	if err != nil {
		return "", fmt.Errorf("failed to compute absolute path for %s (%w)", test.BasePath, err)
	}
	rel, err := relativePath(cwd, abs)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path from %s to %s (%w)", cwd, abs, err)
	}
	return fmt.Sprintf("%s[%s]", rel, test.Test.GetName()), nil
}
