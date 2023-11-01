package discovery

import (
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
)

func DiscoverTests(fileName string, paths ...string) ([]Test, error) {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		return nil, err
	}
	var tests []Test
	for _, folder := range folders {
		test, err := LoadTest(fileName, folder)
		if err != nil {
			return nil, err
		} else if test != nil {
			tests = append(tests, *test)
		}
		// file := filepath.Join(folder, fileName)
		// // check if we have a full test file present
		// if _, err := os.Lstat(file); err == nil {
		// 	apiTests, err := test.Load(file)
		// 	if err == nil {
		// 		basePath := filepath.Dir(file)
		// 		for _, test := range apiTests {
		// 			tests = append(tests, Test{
		// 				Test:     test,
		// 				BasePath: basePath,
		// 			})
		// 		}
		// 	}
		// } else {
		// 	// list files in folder
		// 	files, err := os.ReadDir(folder)
		// 	if err == nil {
		// 		var stepFiles []string
		// 		for _, file := range files {
		// 			fileName := file.Name()
		// 			if !file.IsDir() {
		// 				if StepFileName.MatchString(fileName) {
		// 					stepFiles = append(stepFiles, fileName)
		// 				}
		// 			}
		// 		}
		// 		// we found some steps
		// 		if len(stepFiles) != 0 {
		// 			slices.Sort(stepFiles)
		// 			test := &v1alpha1.Test{
		// 				TypeMeta: metav1.TypeMeta{
		// 					APIVersion: "chainsaw.kyverno.io/v1alpha1",
		// 					Kind:       "Test",
		// 				},
		// 				ObjectMeta: metav1.ObjectMeta{
		// 					Name: filepath.Base(folder),
		// 				},
		// 			}
		// 			stepsMap := map[string]v1alpha1.TestSpecStep{}
		// 			for _, stepFile := range stepFiles {
		// 				groups := StepFileName.FindStringSubmatch(stepFile)
		// 				if steps, err := step.Load(filepath.Join(folder, stepFile)); err != nil {
		// 					fileRef := v1alpha1.FileRef{
		// 						File: stepFile,
		// 					}
		// 					step, ok := stepsMap[groups[1]]
		// 					if !ok {
		// 						step = v1alpha1.TestSpecStep{}
		// 					}
		// 					if step.Name == "" {
		// 						step.Name = groups[2]
		// 					}
		// 					switch groups[2] {
		// 					case "assert":
		// 						step.Spec.Assert = append(step.Spec.Assert, v1alpha1.Assert{
		// 							FileRef: fileRef,
		// 						})
		// 					case "error":
		// 						step.Spec.Error = append(step.Spec.Error, v1alpha1.Error{
		// 							FileRef: fileRef,
		// 						})
		// 					default:
		// 						step.Name = groups[2]
		// 						step.Spec.Apply = append(step.Spec.Apply, v1alpha1.Apply{
		// 							FileRef: fileRef,
		// 						})
		// 					}
		// 					stepsMap[groups[1]] = step
		// 				} else {
		// 					if len(steps) != 1 {
		// 						return nil, fmt.Errorf("more than one test step found in %s", filepath.Join(folder, stepFile))
		// 					}
		// 					stepsMap[groups[1]] = v1alpha1.TestSpecStep{
		// 						Name: steps[0].Name,
		// 						Spec: steps[0].Spec,
		// 					}
		// 				}
		// 			}
		// 			keys := make([]string, 0, len(stepsMap))
		// 			for k := range stepsMap {
		// 				keys = append(keys, k)
		// 			}
		// 			slices.Sort(keys)
		// 			for _, key := range keys {
		// 				test.Spec.Steps = append(test.Spec.Steps, stepsMap[key])
		// 			}
		// 			tests = append(tests, Test{
		// 				Test:     test,
		// 				BasePath: folder,
		// 			})
		// 		}
		// 	}
		// }
	}
	return tests, nil
}
