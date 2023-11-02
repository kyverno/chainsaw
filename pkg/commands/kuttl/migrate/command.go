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
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/resource"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	fileutils "github.com/kyverno/kyverno/ext/file"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
				migrated, err := migrate(out, path, resource)
				if err != nil {
					needsSave = false
					break
				}
				if migrated == nil {
					converted = append(converted, resource)
				} else {
					converted = append(converted, migrated)
					needsSave = true
				}
			}
			if save && needsSave {
				savePath := path
				if !overwrite {
					savePath = strings.TrimRight(path, filepath.Ext(path)) + ".chainsaw.yaml"
				}
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

func migrate(out io.Writer, path string, resource unstructured.Unstructured) (metav1.Object, error) {
	if resource.GetAPIVersion() == "kuttl.dev/v1beta1" {
		switch resource.GetKind() {
		case "TestSuite":
			fmt.Fprintf(out, "Converting %s in %s...\n", "TestSuite", path)
			configuration, err := testSuite(resource)
			if err != nil {
				fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestSuite", path, err)
				return nil, err
			}
			if configuration.GetName() == "" {
				configuration.SetName("configuration")
			}
			return configuration, nil
		case "TestStep":
			fmt.Fprintf(out, "Converting %s in %s...\n", "TestStep", path)
			step, err := testStep(resource)
			if err != nil {
				fmt.Fprintf(out, "  ERROR: failed to convert %s (%s): %s\n", "TestStep", path, err)
				return nil, err
			}
			if step.GetName() == "" {
				groups := discovery.StepFileName.FindStringSubmatch(filepath.Base(path))
				step.SetName(groups[2])
			}
			return step, nil
		case "TestAssert":
			fmt.Fprintf(out, "Converting %s in %s...\n", "TestAssert", path)
			fmt.Fprintf(out, "  ERROR: not supported (%s)\n", path)
			return nil, fmt.Errorf("conversion not supported %s", resource.GetKind())
		default:
			fmt.Fprintf(out, "  ERROR: unknown kuttl resource (%s): %s\n", path, resource.GetKind())
			return nil, fmt.Errorf("unknown kuttl resource %s", resource.GetKind())
		}
	} else {
		return nil, nil
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

	singleOperation := v1alpha1.Operations{}
	for _, operation := range from.Apply {
		singleOperation.Apply = append(singleOperation.Apply, v1alpha1.Apply{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Assert {
		singleOperation.Assert = append(singleOperation.Assert, v1alpha1.Assert{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Error {
		singleOperation.Error = append(singleOperation.Error, v1alpha1.Error{FileRef: v1alpha1.FileRef{File: operation}})
	}
	for _, operation := range from.Delete {
		singleOperation.Delete = append(singleOperation.Delete, v1alpha1.Delete{
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
	to.Spec.Operations = append(to.Spec.Operations, singleOperation)
	return to, nil
}
