package schemas

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "schemas",
		Short:        "Export JSON schemas",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			if _, err := os.Stat(path); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
				if err := os.MkdirAll(path, os.ModeDir|os.ModePerm); err != nil {
					return err
				}
			}
			schemasFs := data.Schemas()
			entries, err := fs.ReadDir(schemasFs, "schemas/json")
			if err != nil {
				return err
			}
			for _, entry := range entries {
				input, err := fs.ReadFile(schemasFs, filepath.Join("schemas/json", entry.Name()))
				if err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(path, entry.Name()), input, os.ModePerm); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}
