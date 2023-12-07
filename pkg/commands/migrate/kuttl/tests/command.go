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
			if err := os.WriteFile(path, data, os.ModePerm); err != nil {
				return err
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
				path := filepath.Join(folder, file)
				fmt.Fprintf(out, "Deleting file %s ...\n", path)
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processStep(out io.Writer, step *v1alpha1.TestSpecStep, s discovery.Step, folder string, save bool) error {
	for f, file := range s.OtherFiles {
		resources, err := resource.Load(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		for i, resource := range resources {
			if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
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
				file := fmt.Sprintf("chainsaw-%s-apply-%d-%d.yaml", step.Name, f+1, i+1)
				if save {
					if err := saveResource(out, folder, file, resource); err != nil {
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
			}
		}
	}
	for f, file := range s.AssertFiles {
		resources, err := resource.Load(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		for i, resource := range resources {
			if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
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
				file := fmt.Sprintf("chainsaw-%s-assert-%d-%d.yaml", step.Name, f+1, i+1)
				if save {
					if err := saveResource(out, folder, file, resource); err != nil {
						return err
					}
				}
				step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
					Assert: &v1alpha1.Assert{
						FileRefOrResource: v1alpha1.FileRefOrResource{
							FileRef: v1alpha1.FileRef{
								File: file,
							},
						},
					},
				})
			}
		}
	}
	for f, file := range s.ErrorFiles {
		resources, err := resource.Load(filepath.Join(folder, file))
		if err != nil {
			return err
		}
		for i, resource := range resources {
			if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
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
				file := fmt.Sprintf("chainsaw-%s-error-%d-%d.yaml", step.Name, f+1, i+1)
				if save {
					if err := saveResource(out, folder, file, resource); err != nil {
						return err
					}
				}
				step.TestStepSpec.Try = append(step.TestStepSpec.Try, v1alpha1.Operation{
					Error: &v1alpha1.Error{
						FileRefOrResource: v1alpha1.FileRefOrResource{
							FileRef: v1alpha1.FileRef{
								File: file,
							},
						},
					},
				})
			}
		}
	}
	return nil
}

func saveResource(out io.Writer, folder, file string, resource unstructured.Unstructured) error {
	path := filepath.Join(folder, file)
	fmt.Fprintf(out, "Saving file %s ...\n", path)
	yamlData, err := yaml.Marshal(&resource)
	if err != nil {
		return fmt.Errorf("converting to yaml: %w", err)
	}
	if err := os.WriteFile(path, yamlData, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func testStep(to *v1alpha1.TestStepSpec, in unstructured.Unstructured) error {
	from, err := convert.To[kuttlapi.TestStep](in)
	// TODO: verify order in kuttl
	if err != nil {
		return err
	}
	for _, operation := range from.Apply {
		to.Try = append(
			to.Try,
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
		to.Try = append(to.Try, v1alpha1.Operation{
			Assert: &v1alpha1.Assert{
				FileRefOrResource: v1alpha1.FileRefOrResource{
					FileRef: v1alpha1.FileRef{
						File: operation,
					},
				},
			},
		})
	}
	for _, operation := range from.Error {
		to.Try = append(to.Try, v1alpha1.Operation{
			Error: &v1alpha1.Error{
				FileRefOrResource: v1alpha1.FileRefOrResource{
					FileRef: v1alpha1.FileRef{
						File: operation,
					},
				},
			},
		})
	}
	for _, operation := range from.Delete {
		to.Try = append(to.Try, v1alpha1.Operation{
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
			return errors.New("found a command with background=true, this is not supported in chainsaw")
		}
		if operation.Namespaced {
			return errors.New("found a command with namespaced=true, this is not supported in chainsaw")
		}
		if operation.IgnoreFailure {
			return errors.New("found a command with ignoreFailure=true, this is not supported in chainsaw")
		}
		if operation.Script != "" {
			to.Try = append(to.Try, v1alpha1.Operation{
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
			to.Try = append(to.Try, v1alpha1.Operation{
				Command: &v1alpha1.Command{
					Timeout:       timeout,
					Entrypoint:    entrypoint,
					Args:          args,
					SkipLogOutput: operation.SkipLogOutput,
				},
			})
		}
	}
	return nil
}

func testAssert(to *v1alpha1.TestStepSpec, in unstructured.Unstructured) error {
	from, err := convert.To[kuttlapi.TestAssert](in)
	if err != nil {
		return err
	}
	for _, cmd := range from.Commands {
		if cmd.Script != "" {
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Script: &v1alpha1.Script{
					Content:       cmd.Script,
					SkipLogOutput: cmd.SkipLogOutput,
				},
			})
		} else if cmd.Command != "" {
			splitCmd, err := shlex.Split(cmd.Command)
			if err != nil {
				return err
			}
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Command: &v1alpha1.Command{
					Entrypoint:    splitCmd[0],
					Args:          splitCmd[1:],
					SkipLogOutput: cmd.SkipLogOutput,
				},
			})
		}
	}
	for _, collector := range from.Collectors {
		switch collector.Type {
		case "pod":
			op := &v1alpha1.PodLogs{
				Name:      collector.Pod,
				Namespace: collector.Namespace,
				Container: collector.Container,
				Selector:  collector.Selector,
			}
			if collector.Tail != 0 {
				op.Tail = ptr.To(collector.Tail)
			}
			to.Catch = append(to.Catch, v1alpha1.Catch{PodLogs: op})
		case "command":
			if collector.Cmd == "" {
				return fmt.Errorf("cmd must be set when tyme is command")
			}
			splitCmd, err := shlex.Split(collector.Cmd)
			if err != nil {
				return err
			}
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Command: &v1alpha1.Command{
					Entrypoint: splitCmd[0],
					Args:       splitCmd[1:],
				},
			})
		case "events":
			to.Catch = append(to.Catch, v1alpha1.Catch{
				Events: &v1alpha1.Events{
					Name:      collector.Pod,
					Namespace: collector.Namespace,
					Selector:  collector.Selector,
				},
			})
		default:
			return fmt.Errorf("unknown collector type: %s", collector.Type)
		}
	}
	return nil
}
