package discovery

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func DiscoverTests2(fileName string, paths ...string) ([]Test, error) {
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
				// check file name pattern
				if !file.IsDir() {
					stepFiles = append(stepFiles, file.Name())
				}
			}
			// we found some steps
			if len(stepFiles) != 0 {
				// TODO sort steps
				test := &v1alpha1.Test{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "chainsaw.kyverno.io/v1alpha1",
						Kind:       "Test",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: filepath.Base(folder),
					},
				}
				for _, stepFile := range stepFiles {
					// TODO: consider case the file contains a test step
					// TODO: if it's not a test step decode operation from file name
					test.Spec.Steps = append(test.Spec.Steps, v1alpha1.TestStepSpec{
						Apply: []v1alpha1.Apply{{
							FileRef: v1alpha1.FileRef{
								File: stepFile,
							},
						}},
					})
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
