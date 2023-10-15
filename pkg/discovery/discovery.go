package discovery

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/test"
	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var stepFileName = regexp.MustCompile(`^(\d\d)-(.*)\.(?:yaml|yml)$`)

type Test struct {
	*v1alpha1.Test
	BasePath string
}

func DiscoverTests(fileName string, paths ...string) ([]Test, error) {
	var files []string
	for _, path := range paths {
		if _, err := os.Lstat(path); err == nil {
			err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Name() == fileName {
					files = append(files, file)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	var tests []Test
	for _, file := range files {
		apiTests, err := test.Load(file)
		if err == nil {
			basePath := filepath.Dir(file)
			for _, test := range apiTests {
				tests = append(tests, Test{
					Test:     test,
					BasePath: basePath,
				})
			}
		}
	}
	return tests, nil
}

func DiscoverTests2(paths ...string) ([]Test, error) {
	var folders []string
	for _, path := range paths {
		if _, err := os.Lstat(path); err == nil {
			err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					folders = append(folders, file)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	var tests []Test
	for _, folder := range folders {
		// list files in folder
		files, err := os.ReadDir(folder)
		if err == nil {
			var stepFiles []string
			for _, file := range files {
				fileName := file.Name()
				if !file.IsDir() {
					if stepFileName.MatchString(fileName) {
						stepFiles = append(stepFiles, fileName)
					}
				}
			}
			// we found some steps
			if len(stepFiles) != 0 {
				slices.Sort(stepFiles)
				test := &v1alpha1.Test{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "chainsaw.kyverno.io/v1alpha1",
						Kind:       "Test",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: filepath.Base(folder),
					},
				}
				stepsMap := map[string]v1alpha1.TestStepSpec{}
				for _, stepFile := range stepFiles {
					// TODO: consider case the file contains a test step
					fileRef := v1alpha1.FileRef{
						File: stepFile,
					}
					groups := stepFileName.FindStringSubmatch(stepFile)
					step, ok := stepsMap[groups[1]]
					if !ok {
						step = v1alpha1.TestStepSpec{}
					}
					switch groups[2] {
					case "assert":
						step.Assert = append(step.Assert, v1alpha1.Assert{
							FileRef: fileRef,
						})
					case "error":
						step.Error = append(step.Error, v1alpha1.Error{
							FileRef: fileRef,
						})
					default:
						step.Apply = append(step.Apply, v1alpha1.Apply{
							FileRef: fileRef,
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
					test.Spec.Steps = append(test.Spec.Steps, stepsMap[key])
				}
				tests = append(tests, Test{
					Test:     test,
					BasePath: folder,
				})
			}
		}
	}
	return tests, nil
}
