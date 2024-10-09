package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const schema = "# yaml-language-server: $schema=https://raw.githubusercontent.com/kyverno/chainsaw/main/.schemas/json/configuration-chainsaw-v1alpha2.json"

func Command() *cobra.Command {
	save := false
	cmd := &cobra.Command{
		Use:          "config",
		Short:        "Upgrade Chainsaw configuration to the latest version",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(
				cmd.OutOrStdout(),
				cmd.ErrOrStderr(),
				save,
				args[0],
			)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, converted files will be saved")
	return cmd
}

func execute(stdout io.Writer, stderr io.Writer, save bool, file string) error {
	c, err := loadConfig(file)
	if err != nil {
		return fmt.Errorf("failed to load configuration file %s: %w", file, err)
	}
	u, err := buildConfigPatch(c)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration patch: %w", err)
	}
	if save {
		fmt.Fprintf(stderr, "Saving file %s ...\n", file)
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
	return nil
}

func loadConfig(file string) (*v1alpha2.Configuration, error) {
	return config.Load(fsutils.NewLocal(), file)
}

func loadDefaultConfig() (*v1alpha2.Configuration, error) {
	return config.DefaultConfiguration()
}

func buildConfigPatch(c *v1alpha2.Configuration) (*unstructured.Unstructured, error) {
	d, err := loadDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}
	o, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default configuration: %w", err)
	}
	m, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal desired configuration: %w", err)
	}
	patch, err := jsonpatch.CreateMergePatch(o, m)
	if err != nil {
		return nil, fmt.Errorf("failed to create patch from default to desired configuration: %w", err)
	}
	var i map[string]any
	if err := json.Unmarshal(patch, &i); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration patch: %w", err)
	}
	u := unstructured.Unstructured{
		Object: i,
	}
	u.SetAPIVersion(v1alpha2.GroupVersion.String())
	u.SetKind("Configuration")
	u.SetName(c.Name)
	u.SetLabels(c.Labels)
	u.SetAnnotations(c.Annotations)
	return &u, nil
}
