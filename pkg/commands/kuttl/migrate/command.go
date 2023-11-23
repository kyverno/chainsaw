package migrate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
	kuttlapi "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	fileutils "github.com/kyverno/kyverno/ext/file"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

func Command() *cobra.Command {
	save := false
	overwrite := false
	cmd := &cobra.Command{
		Use:          "migrate",
		Short:        "Migrate KUTTL tests to Chainsaw",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(cmd.OutOrStdout(), save, overwrite, args...)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, converted files will be saved.")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "If set, overwrites original file.")
	return cmd
}

func execute(out io.Writer, save, overwrite bool, paths ...string) error {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		fmt.Fprintf(out, "  ERROR: failed to discover folders: %s\n", err)
		return err
	}
	for _, folder := range folders {
		if err := processFolder(out, folder, save, overwrite); err != nil {
			fmt.Fprintf(out, "Error processing folder %s: %v\n", folder, err)
		}
	}
	return nil
}

func processFolder(out io.Writer, folder string, save, overwrite bool) error {
	testAssertsMap, err := collectTestAsserts(folder)
	if err != nil {
		fmt.Fprintf(out, "ERROR: failed to collect test asserts: %v\n", err)
		return err
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !fileutils.IsYaml(file.Name()) {
			continue
		}
		path := filepath.Join(folder, file.Name())
		if err := processFile(out, path, save, overwrite, testAssertsMap); err != nil {
			fmt.Fprintf(out, "Error processing file %s: %v\n", path, err)
		}
	}
	return nil
}

func processFile(out io.Writer, path string, save, overwrite bool, testAssertsMap map[string][]v1alpha1.Catch) error {
	resources, err := resource.Load(path)
	if err != nil {
		return err
	}
	var converted []interface{}
	var needsSave bool
	for _, resource := range resources {
		migrated, shouldAppend, err := migrate(out, path, resource, testAssertsMap)
		if err != nil {
			needsSave = false
			break
		}
		if shouldAppend {
			if migrated == nil {
				converted = append(converted, resource)
			} else {
				converted = append(converted, migrated)
				needsSave = true
			}
		}
	}
	if save && needsSave {
		if err := saveConvertedFile(out, path, converted, overwrite); err != nil {
			return err
		}
	}
	return nil
}

func saveConvertedFile(out io.Writer, path string, resources []interface{}, overwrite bool) error {
	savePath := path
	if !overwrite {
		savePath = strings.TrimRight(path, filepath.Ext(path)) + ".chainsaw.yaml"
	}
	fmt.Fprintf(out, "Saving converted file %s to %s...\n", path, savePath)

	var yamlBytes []byte
	for _, res := range resources {
		yamlData, err := yaml.Marshal(res)
		if err != nil {
			return fmt.Errorf("converting to yaml: %w", err)
		}

		yamlBytes = append(yamlBytes, []byte("---\n")...)
		yamlBytes = append(yamlBytes, yamlData...)
	}

	return os.WriteFile(savePath, yamlBytes, os.ModePerm)
}

func migrate(out io.Writer, path string, resource unstructured.Unstructured, testAssertsMap map[string][]v1alpha1.Catch) (interface{}, bool, error) {
	index := extractIndex(path)
	if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
		switch resource.GetKind() {
		case "TestSuite":
			fmt.Fprintf(out, "Converting %s in %s...\n", "TestSuite", path)
			configuration, err := testSuite(resource)
			if err != nil {
				fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestSuite", path, err)
				return nil, false, err
			}
			if configuration.GetName() == "" {
				configuration.SetName("configuration")
			}
			return configuration, true, nil
		case "TestStep":
			groups := discovery.StepFileName.FindStringSubmatch(filepath.Base(path))
			if len(groups) < 3 {
				return nil, true, nil
			}
			fmt.Fprintf(out, "Converting %s in %s...\n", "TestStep", path)
			step, err := testStep(resource)
			if err != nil {
				fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestStep", path, err)
				return nil, false, err
			}
			if step.GetName() == "" {
				step.SetName(strings.ToLower(strings.ReplaceAll(groups[2], "_", "-")))
			}
			// Append TestAsserts (Catch objects) to the TestStep if they exist for this index
			if catchArray, ok := testAssertsMap[index]; ok {
				step.Spec.Catch = append(step.Spec.Catch, catchArray...)
			}
			return step, true, nil
		case "TestAssert":
			return nil, false, nil
		default:
			fmt.Fprintf(out, "  ERROR: unknown kuttl resource (%s): %s\n", path, resource.GetKind())
			return nil, false, fmt.Errorf("unknown kuttl resource %s", resource.GetKind())
		}
	} else {
		return nil, true, nil
	}
}

func testSuite(in unstructured.Unstructured) (*v1alpha1.Configuration, error) {
	from, err := convert.To[kuttlapi.TestSuite](in)
	if err != nil {
		return nil, err
	}
	var timeouts v1alpha1.Timeouts
	if from.Timeout != 0 {
		d := metav1.Duration{Duration: time.Second * time.Duration(from.Timeout)}
		timeouts = v1alpha1.Timeouts{
			Apply:   &d,
			Assert:  &d,
			Error:   &d,
			Delete:  &d,
			Cleanup: &d,
			Exec:    &d,
		}
	}
	to := &v1alpha1.Configuration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "Configuration",
		},
		ObjectMeta: from.ObjectMeta,
		Spec: v1alpha1.ConfigurationSpec{
			Timeouts:     timeouts,
			TestDirs:     from.TestDirs,
			SkipDelete:   from.SkipDelete,
			ReportFormat: v1alpha1.ReportFormatType(from.ReportFormat),
			ReportName:   from.ReportName,
			Namespace:    from.Namespace,
		},
	}
	if from.Parallel != 0 {
		to.Spec.Parallel = &from.Parallel
	}
	return to, nil
}

func testStep(in unstructured.Unstructured) (*v1alpha1.TestStep, error) {
	from, err := convert.To[kuttlapi.TestStep](in)
	// TODO: verify order in kuttl
	if err != nil {
		return nil, err
	}
	to := &v1alpha1.TestStep{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "TestStep",
		},
		ObjectMeta: from.ObjectMeta,
	}
	for _, operation := range from.Apply {
		to.Spec.Try = append(
			to.Spec.Try,
			v1alpha1.Operation{
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: operation,
						},
					},
				},
			},
		)
	}
	for _, operation := range from.Assert {
		to.Spec.Try = append(to.Spec.Try, v1alpha1.Operation{Assert: &v1alpha1.Assert{FileRef: v1alpha1.FileRef{File: operation}}})
	}
	for _, operation := range from.Error {
		to.Spec.Try = append(to.Spec.Try, v1alpha1.Operation{Error: &v1alpha1.Error{FileRef: v1alpha1.FileRef{File: operation}}})
	}
	for _, operation := range from.Delete {
		to.Spec.Try = append(to.Spec.Try, v1alpha1.Operation{
			Delete: &v1alpha1.Delete{
				ObjectReference: v1alpha1.ObjectReference{
					APIVersion: operation.APIVersion,
					Kind:       operation.Kind,
					ObjectSelector: v1alpha1.ObjectSelector{
						Namespace: operation.Namespace,
						Name:      operation.Name,
						Labels:    operation.Labels,
					},
				},
			},
		})
	}
	for _, operation := range from.Commands {
		var timeout *metav1.Duration
		if operation.Timeout != 0 {
			timeout = &metav1.Duration{Duration: time.Second * time.Duration(operation.Timeout)}
		}
		if operation.Background {
			return nil, errors.New("found a command with background=true, this is not supported in chainsaw")
		}
		if operation.Namespaced {
			return nil, errors.New("found a command with namespaced=true, this is not supported in chainsaw")
		}
		if operation.IgnoreFailure {
			return nil, errors.New("found a command with ignoreFailure=true, this is not supported in chainsaw")
		}
		if operation.Script != "" {
			to.Spec.Try = append(to.Spec.Try, v1alpha1.Operation{
				Timeout: timeout,
				Script: &v1alpha1.Script{
					Content:       operation.Script,
					SkipLogOutput: operation.SkipLogOutput,
				},
			})
		} else if operation.Command != "" {
			split, err := shlex.Split(operation.Command)
			if err != nil {
				return nil, err
			}
			entrypoint := split[0]
			var args []string
			if len(split) > 1 {
				args = split[1:]
			}
			to.Spec.Try = append(to.Spec.Try, v1alpha1.Operation{
				Timeout: timeout,
				Command: &v1alpha1.Command{
					Entrypoint:    entrypoint,
					Args:          args,
					SkipLogOutput: operation.SkipLogOutput,
				},
			})
		}
	}
	return to, nil
}

// MigrateTestAssert migrates a KUTTL TestAssert to a Chainsaw TestStep.
func migrateTestAssert(in unstructured.Unstructured) ([]v1alpha1.Catch, error) {
	from, err := convert.To[kuttlapi.TestAssert](in)
	if err != nil {
		return nil, err
	}
	var catchArray []v1alpha1.Catch

	// Handle TestAssertCommands
	for _, cmd := range from.Commands {
		var catch v1alpha1.Catch
		if cmd.Script != "" {
			catch = v1alpha1.Catch{
				Script: &v1alpha1.Script{
					Content:       cmd.Script,
					SkipLogOutput: cmd.SkipLogOutput,
				},
			}

		} else if cmd.Command != "" {
			splitCmd, err := shlex.Split(cmd.Command)
			if err != nil {
				return nil, err
			}
			catch = v1alpha1.Catch{
				Command: &v1alpha1.Command{
					Entrypoint:    splitCmd[0],
					Args:          splitCmd[1:],
					SkipLogOutput: cmd.SkipLogOutput,
				},
			}
		}
		catchArray = append(catchArray, catch)
	}

	// Handle Collectors
	for _, collector := range from.Collectors {
		var catch v1alpha1.Catch
		switch collector.Type {
		case "pod":
			catch = v1alpha1.Catch{
				PodLogs: &v1alpha1.PodLogs{
					Name:      collector.Pod,
					Namespace: collector.Namespace,
					Container: collector.Container,
					Selector:  collector.Selector,
					Tail:      ptr.To(collector.Tail),
				},
			}
		case "command":
			if collector.Cmd != "" {
				splitCmd, err := shlex.Split(collector.Cmd)
				if err != nil {
					return nil, err
				}
				catch = v1alpha1.Catch{
					Command: &v1alpha1.Command{
						Entrypoint: splitCmd[0],
						Args:       splitCmd[1:],
					},
				}
			}
		case "events":
			catch = v1alpha1.Catch{
				Events: &v1alpha1.Events{
					Name:      collector.Pod,
					Namespace: collector.Namespace,
					Selector:  collector.Selector,
				},
			}
		default:
			return nil, fmt.Errorf("unknown collector type: %s", collector.Type)
		}
		catchArray = append(catchArray, catch)
	}
	return catchArray, nil
}

func extractIndex(path string) string {
	baseName := filepath.Base(path)
	splitName := strings.Split(baseName, "-")
	if len(splitName) > 0 {
		return splitName[0]
	}
	return ""
}

func collectTestAsserts(folder string) (map[string][]v1alpha1.Catch, error) {
	testAssertsMap := make(map[string][]v1alpha1.Catch)
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !fileutils.IsYaml(file.Name()) {
			continue
		}
		path := filepath.Join(folder, file.Name())
		resources, err := resource.Load(path)
		if err != nil {
			return nil, err
		}

		for _, res := range resources {
			if res.GetAPIVersion() == "kuttl.dev/v1beta1" && res.GetKind() == "TestAssert" {
				index := extractIndex(path)
				catchArray, err := migrateTestAssert(res)
				if err != nil {
					return nil, err
				}
				testAssertsMap[index] = append(testAssertsMap[index], catchArray...)
			}
		}
	}
	return testAssertsMap, nil
}
