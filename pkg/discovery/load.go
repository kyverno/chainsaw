package discovery

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/step"
	"github.com/kyverno/chainsaw/pkg/test"
	"go.uber.org/multierr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tryLoadTest(file string) (*v1alpha1.Test, error) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	tests, err := test.Load(file)
	if err != nil {
		return nil, err
	}
	if len(tests) != 1 {
		return nil, fmt.Errorf("found more than one test in %s (%d)", file, len(tests))
	}
	return tests[0], nil
}

func TryFindStepFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	} else {
		// collect and sort candidate files
		var stepFiles []string
		for _, file := range files {
			fileName := file.Name()
			if !file.IsDir() {
				if StepFileName.MatchString(fileName) {
					stepFiles = append(stepFiles, fileName)
				}
			}
		}
		if len(stepFiles) != 0 {
			slices.Sort(stepFiles)
		}
		return stepFiles, nil
	}
}

func LoadTest(fileName string, path string) (*Test, error) {
	// first, try to load a test manifest
	var errs []error
	var apiTest *v1alpha1.Test
	if path == "" {
		return nil, errors.New("path must be specified")
	}
	if fileName != "" {
		if test, err := tryLoadTest(filepath.Join(path, fileName)); err != nil {
			errs = append(errs, fmt.Errorf("failed to load test file (%w)", err))
		} else if test != nil {
			apiTest = test
		}
	}
	if apiTest == nil {
		apiTest = &v1alpha1.Test{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: filepath.Base(path),
			},
		}
	}
	// next, look at files
	files, err := TryFindStepFiles(path)
	if err != nil {
		return nil, err
	} else {
		if len(files) > 0 && len(apiTest.Spec.Steps) != 0 {
			errs = append(errs, errors.New("test has steps and files matched the convention, files will be ignored"))
		} else {
			// load test step resources first
			var manifestFiles []string
			stepsMap := map[string]v1alpha1.TestSpecStep{}
			for _, file := range files {
				if steps, err := step.Load(filepath.Join(path, file)); err != nil {
					// errs = append(errs, fmt.Errorf("failed to load test step file (%w)", err))
					manifestFiles = append(manifestFiles, file)
				} else if len(steps) != 1 {
					errs = append(errs, fmt.Errorf("more than one test step found in %s", filepath.Join(path, file)))
				} else {
					groups := StepFileName.FindStringSubmatch(file)
					stepsMap[groups[1]] = v1alpha1.TestSpecStep{
						Name:         steps[0].Name,
						TestStepSpec: steps[0].Spec,
					}
				}
			}
			var assertFiles []string
			var errorFiles []string
			var otherFiles []string
			for _, file := range manifestFiles {
				groups := StepFileName.FindStringSubmatch(file)
				switch groups[2] {
				case "assert":
					assertFiles = append(assertFiles, file)
				case "errors":
					errorFiles = append(errorFiles, file)
				default:
					otherFiles = append(otherFiles, file)
				}
			}
			manifestFiles = []string{}
			manifestFiles = append(manifestFiles, otherFiles...)
			manifestFiles = append(manifestFiles, assertFiles...)
			manifestFiles = append(manifestFiles, errorFiles...)
			for _, file := range manifestFiles {
				groups := StepFileName.FindStringSubmatch(file)
				step, ok := stepsMap[groups[1]]
				if !ok {
					step = v1alpha1.TestSpecStep{}
				}
				if step.Name == "" {
					step.Name = groups[2]
				}
				fileRef := v1alpha1.FileRef{
					File: file,
				}
				switch groups[2] {
				case "assert":
					step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
						Assert: &v1alpha1.Assert{
							FileRefOrResource: v1alpha1.FileRefOrResource{
								FileRef: fileRef,
							},
						},
					})
				case "errors":
					step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
						Error: &v1alpha1.Error{
							FileRefOrResource: v1alpha1.FileRefOrResource{
								FileRef: fileRef,
							},
						},
					})
				default:
					step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
						Apply: &v1alpha1.Apply{
							FileRefOrResource: v1alpha1.FileRefOrResource{
								FileRef: fileRef,
							},
						},
					})
				}
				stepsMap[groups[1]] = step
			}
			keys := make([]string, 0, len(stepsMap))
			for k := range stepsMap {
				keys = append(keys, k)
			}
			slices.Sort(keys)
			for _, key := range keys {
				apiTest.Spec.Steps = append(apiTest.Spec.Steps, stepsMap[key])
			}
		}
	}
	if len(apiTest.Spec.Steps) == 0 && errs == nil {
		return nil, nil
	}
	return &Test{
		Test:     apiTest,
		BasePath: path,
		Err:      multierr.Combine(errs...),
	}, nil
}
