package discovery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/test"
	"golang.org/x/exp/maps"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tryLoadTestFile(file string, remarshal bool) ([]*v1alpha1.Test, error) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	tests, err := test.Load(file, remarshal)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func tryLoadTestFiles(fileName string, path string, remarshal bool) ([]*v1alpha1.Test, error) {
	if filepath.Ext(fileName) != "" {
		return tryLoadTestFile(filepath.Join(path, fileName), remarshal)
	}
	tests, err := tryLoadTestFile(filepath.Join(path, fileName+".yaml"), remarshal)
	if err != nil {
		return nil, err
	}
	if tests != nil {
		return tests, nil
	}
	return tryLoadTestFile(filepath.Join(path, fileName+".yml"), remarshal)
}

func LoadTest(fileName string, path string, remarshal bool) ([]Test, error) {
	// first, try to load a test manifest
	if path == "" {
		return nil, errors.New("path must be specified")
	}
	var tests []Test
	if fileName != "" {
		apiTests, err := tryLoadTestFiles(fileName, path, remarshal)
		if err != nil {
			return nil, err
		}
		if len(apiTests) != 0 {
			for _, apiTest := range apiTests {
				tests = append(tests, Test{
					Test:     apiTest,
					BasePath: path,
					Err:      nil,
				})
			}
			return tests, nil
		}
	}
	// next, look at files
	steps, err := TryFindStepFiles(path)
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, nil
	}
	keys := maps.Keys(steps)
	slices.Sort(keys)
	test := &v1alpha1.Test{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "Test",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(strings.ReplaceAll(filepath.Base(path), "_", "-")),
		},
	}
	for _, key := range keys {
		step := v1alpha1.TestStep{
			Name: fmt.Sprintf("step-%s", key),
		}
		for _, file := range steps[key].OtherFiles {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Apply: &v1alpha1.Apply{
					ActionResourceRef: v1alpha1.ActionResourceRef{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
		}
		for _, file := range steps[key].AssertFiles {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Assert: &v1alpha1.Assert{
					ActionCheckRef: v1alpha1.ActionCheckRef{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
		}
		for _, file := range steps[key].ErrorFiles {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Error: &v1alpha1.Error{
					ActionCheckRef: v1alpha1.ActionCheckRef{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
		}
		test.Spec.Steps = append(test.Spec.Steps, step)
	}
	tests = append(tests, Test{
		Test:     test,
		BasePath: path,
		Err:      nil,
	})
	return tests, nil
}
