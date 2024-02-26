package tests

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/shlex"
	kuttlapi "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

func Command() *cobra.Command {
	save := false
	cleanup := false
	cmd := &cobra.Command{
		Use:          "tests",
		Short:        "Migrate KUTTL tests to Chainsaw",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(cmd.OutOrStdout(), save, cleanup, args...)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, converted files will be saved")
	cmd.Flags().BoolVar(&cleanup, "cleanup", false, "If set, delete converted files")
	return cmd
}

func execute(out io.Writer, save, cleanup bool, paths ...string) error {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		fmt.Fprintf(out, "ERROR: failed to discover folders: %s\n", err)
		return err
	}
	for _, folder := range folders {
		if err := processFolder(out, folder, save, cleanup); err != nil {
			fmt.Fprintf(out, "ERROR: failed to process folder %s: %v\n", folder, err)
		}
	}
	return nil
}

func processFolder(out io.Writer, folder string, save, cleanup bool) error {
	if _, err := os.Stat(filepath.Join(folder, "chainsaw-test.yaml")); err == nil {
		return nil
	}
	steps, err := discovery.TryFindStepFiles(folder)
	if err != nil {
		fmt.Fprintf(out, "ERROR: failed to collect test files: %v\n", err)
		return err
	}
	if len(steps) != 0 {
		fmt.Fprintf(out, "Converting test %s ...\n", folder)
		keys := maps.Keys(steps)
		slices.Sort(keys)
		test := v1alpha1.Test{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.SchemeGroupVersion.String(),
				Kind:       "Test",
			},
		}
		test.SetName(strings.ToLower(strings.ReplaceAll(filepath.Base(folder), "_", "-")))
		for _, key := range keys {
			step := v1alpha1.TestSpecStep{
				Name: fmt.Sprintf("step-%s", key),
			}
			if err := processStep(out, &step, steps[key], folder, save); err != nil {
				return err
			}
			test.Spec.Steps = append(test.Spec.Steps, step)
		}
		data, err := yaml.Marshal(&test)
		if err != nil {
			return err
		}
		if save {
			path := filepath.Join(folder, "chainsaw-test.yaml")
			fmt.Fprintf(out, "Saving file %s ...\n", path)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()
			if _, err := f.WriteString("# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/test-chainsaw-v1alpha1.json\n"); err != nil {
				return fmt.Errorf("failed to write in file: %w", err)
			}
			if _, err := f.Write(data); err != nil {
				return fmt.Errorf("failed to write in file: %w", err)
			}
		} else {
			fmt.Fprintln(out, string(data))
		}
		if save && cleanup {
			var files []string
			for _, step := range steps {
				files = append(files, step.AssertFiles...)
				files = append(files, step.ErrorFiles...)
				files = append(files, step.OtherFiles...)
			}
			for _, file := range files {
				if file != "" {
					path := filepath.Join(folder, file)
					fmt.Fprintf(out, "Deleting file %s ...\n", path)
					if err := os.Remove(path); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func isKuttl(resource unstructured.Unstructured) bool {
	return strings.HasPrefix(resource.GetAPIVersion(), "kuttl.dev/")
}

func processStep(out io.Writer, step *v1alpha1.TestSpecStep, s discovery.Step, folder string, save bool) error {
	for f, file := range s.OtherFiles {
		resources, err := resource.Load(filepath.Join(folder, file), true)
		if err != nil {
			return err
		}
		containsKuttlResources := false
		for _, resource := range resources {
			if isKuttl(resource) {
				containsKuttlResources = true
			}
		}
		if !containsKuttlResources {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.OtherFiles[f] = ""
			continue
		}
		var filteredResources []unstructured.Unstructured
		for _, resource := range resources {
			if isKuttl(resource) {
				switch resource.GetKind() {
				case "TestStep":
					err := testStep(&step.TestStepSpec, resource)
					if err != nil {
						fmt.Fprintf(out, "ERROR: failed to convert %s (%s): %s\n", "TestStep", filepath.Join(folder, file), err)
						return err
					}
				case "TestAssert":
					err := testAssert(&step.TestStepSpec, resource)
					if err != nil {
						fmt.Fprintf(out, "ERROR: failed to convert %s (%s): %s\n", "TestAssert", filepath.Join(folder, file), err)
						return err
					}
				default:
					return fmt.Errorf("type not supported %s / %s", resource.GetAPIVersion(), resource.GetKind())
				}
			} else {
				filteredResources = append(filteredResources, resource)
			}
		}
		if len(filteredResources) != 0 {
			if save {
				if err := saveResources(out, folder, file, filteredResources...); err != nil {
					return err
				}
			}
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.OtherFiles[f] = ""
		}
	}
	for f, file := range s.AssertFiles {
		resources, err := resource.Load(filepath.Join(folder, file), true)
		if err != nil {
			return err
		}
		containsKuttlResources := false
		for _, resource := range resources {
			if isKuttl(resource) {
				containsKuttlResources = true
			}
		}
		if !containsKuttlResources {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Assert: &v1alpha1.Assert{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.AssertFiles[f] = ""
			continue
		}
		var filteredResources []unstructured.Unstructured
		for _, resource := range resources {
			if isKuttl(resource) {
				switch resource.GetKind() {
				case "TestAssert":
					err := testAssert(&step.TestStepSpec, resource)
					if err != nil {
						fmt.Fprintf(out, "ERROR: failed to convert %s (%s): %s\n", "TestAssert", filepath.Join(folder, file), err)
						return err
					}
				default:
					return fmt.Errorf("type not supported %s / %s", resource.GetAPIVersion(), resource.GetKind())
				}
			} else {
				filteredResources = append(filteredResources, resource)
			}
		}
		if len(filteredResources) != 0 {
			if save {
				if err := saveResources(out, folder, file, filteredResources...); err != nil {
					return err
				}
			}
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Assert: &v1alpha1.Assert{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.AssertFiles[f] = ""
		}
	}
	for f, file := range s.ErrorFiles {
		resources, err := resource.Load(filepath.Join(folder, file), true)
		if err != nil {
			return err
		}
		containsKuttlResources := false
		for _, resource := range resources {
			if isKuttl(resource) {
				containsKuttlResources = true
			}
		}
		if !containsKuttlResources {
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Error: &v1alpha1.Error{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.ErrorFiles[f] = ""
			continue
		}
		var filteredResources []unstructured.Unstructured
		for _, resource := range resources {
			if isKuttl(resource) {
				switch resource.GetKind() {
				case "TestAssert":
					err := testAssert(&step.TestStepSpec, resource)
					if err != nil {
						fmt.Fprintf(out, "ERROR: failed to convert %s (%s): %s\n", "TestAssert", filepath.Join(folder, file), err)
						return err
					}
				default:
					return fmt.Errorf("type not supported %s / %s", resource.GetAPIVersion(), resource.GetKind())
				}
			} else {
				filteredResources = append(filteredResources, resource)
			}
		}
		if len(filteredResources) != 0 {
			if save {
				if err := saveResources(out, folder, file, filteredResources...); err != nil {
					return err
				}
			}
			step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
				Error: &v1alpha1.Error{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: file,
						},
					},
				},
			})
			// no cleanup
			s.ErrorFiles[f] = ""
		}
	}
	return nil
}

func saveResources(out io.Writer, folder, file string, resources ...unstructured.Unstructured) error {
	path := filepath.Join(folder, file)
	fmt.Fprintf(out, "Saving file %s ...\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	for i := range resources {
		yamlData, err := yaml.Marshal(&resources[i])
		if err != nil {
			return fmt.Errorf("converting to yaml: %w", err)
		}
		if _, err := f.Write(yamlData); err != nil {
			return fmt.Errorf("failed to write in file: %w", err)
		}
		if i < len(resources)-1 {
			if _, err := f.WriteString("---\n"); err != nil {
				return fmt.Errorf("failed to write in file: %w", err)
			}
		}
	}
	return nil
}

func prepend[T any](slice []T, elems ...T) []T {
	var out []T
	if len(elems) != 0 {
		out = append(out, elems...)
	}
	if len(slice) != 0 {
		out = append(out, slice...)
	}
	return out
}

func testStep(to *v1alpha1.TestStepSpec, in unstructured.Unstructured) error {
	from, err := convert.To[kuttlapi.TestStep](in)
	if err != nil {
		return err
	}
	var operations []v1alpha1.Operation
	for _, operation := range from.Commands {
		var timeout *metav1.Duration
		if operation.Timeout != 0 {
			timeout = &metav1.Duration{Duration: time.Second * time.Duration(operation.Timeout)}
		}
		if operation.Background {
			return errors.New("found a command with background=true, this is not supported in chainsaw")
		}
		if operation.Namespaced {
			return errors.New("found a command with namespaced=true, this is not supported in chainsaw")
		}
		if operation.IgnoreFailure {
			return errors.New("found a command with ignoreFailure=true, this is not supported in chainsaw")
		}
		if operation.Script != "" {
			operations = append(operations, v1alpha1.Operation{
				Script: &v1alpha1.Script{
					Timeout:       timeout,
					Content:       operation.Script,
					SkipLogOutput: operation.SkipLogOutput,
				},
			})
		} else if operation.Command != "" {
			split, err := shlex.Split(operation.Command)
			if err != nil {
				return err
			}
			entrypoint := split[0]
			var args []string
			if len(split) > 1 {
				args = split[1:]
			}
			operations = append(operations, v1alpha1.Operation{
				Command: &v1alpha1.Command{
					Timeout:       timeout,
					Entrypoint:    entrypoint,
					Args:          args,
					SkipLogOutput: operation.SkipLogOutput,
				},
			})
		}
	}
	for _, operation := range from.Apply {
		operations = append(operations, v1alpha1.Operation{
			Apply: &v1alpha1.Apply{
				FileRefOrResource: v1alpha1.FileRefOrResource{
					FileRef: v1alpha1.FileRef{
						File: operation,
					},
				},
			},
		})
	}
	for _, operation := range from.Assert {
		operations = append(operations, v1alpha1.Operation{
			Assert: &v1alpha1.Assert{
				FileRefOrCheck: v1alpha1.FileRefOrCheck{
					FileRef: v1alpha1.FileRef{
						File: operation,
					},
				},
			},
		})
	}
	for _, operation := range from.Error {
		operations = append(operations, v1alpha1.Operation{
			Error: &v1alpha1.Error{
				FileRefOrCheck: v1alpha1.FileRefOrCheck{
					FileRef: v1alpha1.FileRef{
						File: operation,
					},
				},
			},
		})
	}
	for _, operation := range from.Delete {
		operations = append(operations, v1alpha1.Operation{
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
	to.Try = prepend(to.Try, operations...)
	return nil
}

func testAssert(to *v1alpha1.TestStepSpec, in unstructured.Unstructured) error {
	from, err := convert.To[kuttlapi.TestAssert](in)
	if err != nil {
		return err
	}
	// TODO: timeout
	for _, cmd := range from.Commands {
		if cmd.Script != "" {
			to.Try = append(to.Try, v1alpha1.Operation{
				Script: &v1alpha1.Script{
					Content:       cmd.Script,
					SkipLogOutput: cmd.SkipLogOutput,
				},
			})
		} else if cmd.Command != "" {
			to.Try = append(to.Try, v1alpha1.Operation{
				Script: &v1alpha1.Script{
					Content:       cmd.Command,
					SkipLogOutput: cmd.SkipLogOutput,
				},
			})
		}
	}
	for _, collector := range from.Collectors {
		if collector.Type == "" && collector.Cmd != "" {
			collector.Type = "command"
		}
		if collector.Type == "" && (collector.Pod != "" || collector.Selector != "") {
			collector.Type = "pod"
		}
		switch collector.Type {
		case "pod":
			op := &v1alpha1.PodLogs{
				ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
					Name:      collector.Pod,
					Namespace: collector.Namespace,
					Selector:  collector.Selector,
				},
				Container: collector.Container,
			}
			if collector.Tail != 0 {
				op.Tail = ptr.To(collector.Tail)
			}
			to.Catch = append(to.Catch, v1alpha1.Catch{PodLogs: op})
		case "command":
			if collector.Cmd == "" {
				return fmt.Errorf("cmd must be set when tyme is command")
			}
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Script: &v1alpha1.Script{
					Content: collector.Cmd,
				},
			})
		case "events":
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Events: &v1alpha1.Events{
					ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
						Name:      collector.Pod,
						Namespace: collector.Namespace,
						Selector:  collector.Selector,
					},
				},
			})
		default:
			return fmt.Errorf("unknown collector type: %s", collector.Type)
		}
	}
	return nil
}
