package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	kuttlapi "github.com/kudobuilder/kuttl/pkg/apis/testharness/v1beta1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/loaders/resource"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const schema = "# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/configuration-chainsaw-v1alpha1.json"

func Command() *cobra.Command {
	save := false
	cleanup := false
	cmd := &cobra.Command{
		Use:          "config",
		Short:        "Migrate KUTTL config to Chainsaw",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(
				cmd.OutOrStdout(),
				cmd.ErrOrStderr(),
				save,
				cleanup,
				args[0],
			)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, converted files will be saved")
	cmd.Flags().BoolVar(&cleanup, "cleanup", false, "If set, delete converted files")
	return cmd
}

func execute(stdout io.Writer, stderr io.Writer, save, cleanup bool, path string) error {
	resources, err := resource.Load(path, true)
	if err != nil {
		return err
	}
	if len(resources) != 1 {
		return fmt.Errorf("invalid number of resources found (%d) in %s", len(resources), path)
	}
	cfg, err := migrate(stderr, path, resources[0])
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("converting to yaml: %w", err)
	}
	if save {
		path := filepath.Join(filepath.Dir(path), config.DefaultFileName)
		fmt.Fprintf(stderr, "Saving file %s ...\n", path)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		if _, err := f.WriteString(schema + "\n"); err != nil {
			return fmt.Errorf("failed to write in file: %w", err)
		}
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("failed to write in file: %w", err)
		}
	} else {
		fmt.Fprintln(stdout, schema)
		fmt.Fprintln(stdout, string(data))
	}
	if save && cleanup {
		fmt.Fprintf(stderr, "Deleting file %s ...\n", path)
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func isKuttl(resource unstructured.Unstructured) bool {
	return strings.HasPrefix(resource.GetAPIVersion(), "kuttl.dev/")
}

func migrate(stderr io.Writer, path string, resource unstructured.Unstructured) (*v1alpha1.Configuration, error) {
	fmt.Fprintf(stderr, "Converting config %s ...\n", path)
	if isKuttl(resource) {
		switch resource.GetKind() {
		case "TestSuite":
			configuration, err := testSuite(resource)
			if err != nil {
				fmt.Fprintf(stderr, "ERROR: failed to convert %s (%s): %s\n", "TestSuite", path, err)
				return nil, err
			}
			if configuration.GetName() == "" {
				configuration.SetName("configuration")
			}
			return configuration, nil
		default:
			fmt.Fprintf(stderr, "ERROR: unknown kuttl resource (%s): %s\n", path, resource.GetKind())
			return nil, fmt.Errorf("unknown kuttl resource %s", resource.GetKind())
		}
	}
	return nil, fmt.Errorf("unknown resource %s / %s", resource.GetAPIVersion(), resource.GetKind())
}

func testSuite(in unstructured.Unstructured) (*v1alpha1.Configuration, error) {
	from, err := convert.To[kuttlapi.TestSuite](in)
	if err != nil {
		return nil, err
	}
	var timeouts v1alpha1.DefaultTimeouts
	if from.Timeout != 0 {
		d := metav1.Duration{Duration: time.Second * time.Duration(from.Timeout)}
		timeouts = v1alpha1.DefaultTimeouts{
			Apply:   d,
			Assert:  d,
			Error:   d,
			Delete:  d,
			Cleanup: d,
			Exec:    d,
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
