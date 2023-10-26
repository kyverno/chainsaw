package migrate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	kuttlapi "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/utils/convert"
	fileutils "github.com/kyverno/chainsaw/pkg/utils/file"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func Command() *cobra.Command {
	save := false
	cmd := &cobra.Command{
		Use:          "migrate",
		Short:        "Migrate KUTTL tests to Chainsaw",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(cmd.OutOrStdout(), save, args...)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, converted files will be saved.")
	return cmd
}

func execute(out io.Writer, save bool, paths ...string) error {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		fmt.Fprintf(out, "  ERROR: failed to discover folders: %s\n", err)
		return err
	}
	for _, folder := range folders {
		files, err := os.ReadDir(folder)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if !fileutils.IsYaml(fileName) {
				continue
			}
			path := filepath.Join(folder, fileName)
			resources, err := resource.Load(path)
			if err != nil {
				continue
			}
			var converted []interface{}
			needsSave := false
			for _, resource := range resources {
				if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
					switch resource.GetKind() {
					case "TestSuite":
						fmt.Fprintf(out, "Converting %s in %s...\n", "TestSuite", path)
						configuration, err := testSuite(resource)
						if err != nil {
							fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestSuite", path, err)
							return err
						}
						needsSave = true
						converted = append(converted, configuration)
					case "TestStep":
						fmt.Fprintf(out, "Converting %s in %s...\n", "TestStep", path)
						step, err := testStep(resource)
						if err != nil {
							fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestStep", path, err)
							return err
						}
						needsSave = true
						converted = append(converted, step)
					case "TestAssert":
						fmt.Fprintf(out, "Converting %s in %s...\n", "TestAssert", path)
						fmt.Fprintf(out, "  ERROR: not supported (%s)\n", path)
						return fmt.Errorf("conversion not supported %s", resource.GetKind())
					default:
						fmt.Fprintf(out, "  ERROR: unknown kuttl resource (%s): %s\n", path, err)
						return fmt.Errorf("unknown kuttl resource %s", resource.GetKind())
					}
				} else {
					converted = append(converted, resource)
				}
			}
			if save && needsSave {
				savePath := strings.TrimRight(path, filepath.Ext(path)) + ".chainsaw.yaml"
				fmt.Fprintf(out, "Saving converted file %s to %s...\n", path, savePath)
				var yamlBytes []byte
				for _, resource := range converted {
					finalBytes, err := yaml.Marshal(resource)
					if err != nil {
						fmt.Fprintf(out, "  ERROR: converting to yaml: %s\n", err)
						return err
					}
					yamlBytes = append(yamlBytes, []byte("---\n")...)
					yamlBytes = append(yamlBytes, finalBytes...)
				}
				if err := os.WriteFile(savePath, yamlBytes, os.ModePerm); err != nil {
					fmt.Fprintf(out, "  ERROR: saving file (%s): %s\n", savePath, err)
					return err
				}
			}
		}
	}
	return nil
}

func testSuite(in unstructured.Unstructured) (*v1alpha1.Configuration, error) {
	from, err := convert.To[kuttlapi.TestSuite](in)
	if err != nil {
		return nil, err
	}
	timeout := &metav1.Duration{
		Duration: time.Second * 30,
	}
	if from.Timeout != 0 {
		timeout.Duration = time.Second * time.Duration(from.Timeout)
	}
	to := &v1alpha1.Configuration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "Configuration",
		},
		ObjectMeta: from.ObjectMeta,
		Spec: v1alpha1.ConfigurationSpec{
			Timeout:      timeout,
			TestDirs:     from.TestDirs,
			SkipDelete:   from.SkipDelete,
			Parallel:     from.Parallel,
			ReportFormat: v1alpha1.ReportFormatType(from.ReportFormat),
			ReportName:   from.ReportName,
			Namespace:    from.Namespace,
		},
	}
	return to, nil
}

func testStep(in unstructured.Unstructured) (*v1alpha1.TestStep, error) {
	from, err := convert.To[kuttlapi.TestStep](in)
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
		to.Spec.Apply = append(to.Spec.Apply, v1alpha1.Apply{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Assert {
		to.Spec.Assert = append(to.Spec.Assert, v1alpha1.Assert{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Error {
		to.Spec.Error = append(to.Spec.Error, v1alpha1.Error{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Delete {
		to.Spec.Delete = append(to.Spec.Delete, v1alpha1.Delete{
			ObjectReference: v1alpha1.ObjectReference{
				APIVersion: operation.APIVersion,
				Kind:       operation.Kind,
				ObjectSelector: v1alpha1.ObjectSelector{
					Namespace: operation.Namespace,
					Name:      operation.Name,
					Labels:    operation.Labels,
				},
			},
		})
	}
	return to, nil
}
