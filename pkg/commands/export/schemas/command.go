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
			outPath := args[0]
			if _, err := os.Stat(outPath); err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
				if err := os.MkdirAll(outPath, os.ModeDir|os.ModePerm); err != nil {
					return err
				}
			}
			schemasFs, err := data.Schemas()
			if err != nil {
				return err
			}
			entries, err := fs.ReadDir(schemasFs, ".")
			if err != nil {
				return err
			}
			for _, entry := range entries {
				input, err := fs.ReadFile(schemasFs, entry.Name())
				if err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(outPath, entry.Name()), input, 0o600); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}
